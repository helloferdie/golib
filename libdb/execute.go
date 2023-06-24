package libdb

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/helloferdie/golib/liblogger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// returningID - return ID upon query exection, only applicable for postgres
type returningID struct {
	ID int64 `db:"id"`
}

// Exec - Execute query
func Exec(d *sqlx.DB, query string, values map[string]interface{}) (int64, int64, error) {
	result, err := d.NamedExec(query, values)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error execute query %v", err)
		return 0, 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	return id, rows, nil
}

// Get - Get single row from query
func Get(d *sqlx.DB, list interface{}, query string, values map[string]interface{}) (bool, error) {
	exist := false
	rows, err := d.NamedQuery(query, values)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error get query %v", err)
		return exist, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(list)
		if err != nil {
			liblogger.Log(nil, true).Errorf("Error scan row %v", err)
			return exist, err
		}
		exist = true
	}
	rows.Close()
	return exist, nil
}

// GetByField - Get single row based on provided fields from query
func GetByField(d *sqlx.DB, cfg Config, dt interface{}, params map[string]interface{}, condition string) (bool, error) {
	exist, err := Get(d, dt, "SELECT "+cfg.Fields+" FROM "+cfg.Table+" WHERE 1=1 "+condition, params)
	return exist, err
}

// GetByID - Get single row by ID from query
func GetByID(d *sqlx.DB, cfg Config, dt interface{}, id interface{}) (bool, error) {
	exist, err := GetByField(d, cfg, dt, map[string]interface{}{
		"id": id,
	}, "AND id = :id "+cfg.GetConditionSoftDelete())
	return exist, err
}

// GetByUUID - Get single row by UUID from query
func GetByUUID(d *sqlx.DB, cfg Config, dt interface{}, uuid string) (bool, error) {
	exist, err := GetByField(d, cfg, dt, map[string]interface{}{
		"uuid": uuid,
	}, "AND uuid = :uuid "+cfg.GetConditionSoftDelete())
	return exist, err
}

// GetSoftDeleteByID - Get soft deleted row by ID from query
func GetSoftDeleteByID(d *sqlx.DB, cfg Config, dt interface{}, id int64) (bool, error) {
	exist, err := GetByField(d, cfg, dt, map[string]interface{}{
		"id": id,
	}, "AND id = :id AND deleted_at IS NOT NULL ")
	return exist, err
}

// Select - Select rows from query
func Select(d *sqlx.DB, list interface{}, query string, values map[string]interface{}) error {
	nstmt, err := d.PrepareNamed(query)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error select prepare named query %v", err)
		return err
	}
	err = nstmt.Select(list, values)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error select query %v", err)
		return err
	}
	return nil
}

// List - Get slices of return data from query
func List(d *sqlx.DB, cfg Config, list interface{}, conditionVal map[string]interface{}, condition string, pagination *ModelPaginationRequest) (int64, error) {
	totalItems, err := ListByField(d, list, conditionVal, cfg.GetConditionSoftDelete()+condition, cfg.Table, cfg.Table+".id", cfg.Fields, pagination)
	return totalItems, err
}

// ListByField - Get slices of return data from query
func ListByField(d *sqlx.DB, list interface{}, conditionVal map[string]interface{}, condition string, table string, fieldCount string, fields string, pagination *ModelPaginationRequest) (int64, error) {
	t := new(ModelTotal)
	_, err := Get(d, t, "SELECT COUNT("+fieldCount+") AS total FROM "+table+" WHERE 1=1 "+condition, conditionVal)
	if err != nil {
		fmt.Println(err)
		return t.Total, err
	}

	orderQuery, orderValues := PrepareOrderQuery(pagination)
	for k, v := range orderValues {
		conditionVal[k] = v
	}

	nstmt, err := d.PrepareNamed("SELECT " + fields + " FROM " + table + " WHERE 1=1 " + condition + orderQuery)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error select prepare named query %v", err)
		return t.Total, err
	}

	err = nstmt.Select(list, conditionVal)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error select query %v", err)
		return t.Total, err
	}
	return t.Total, nil
}

// ListRaw - Raw query list
func ListRaw(d *sqlx.DB, list interface{}, query string, conditionVal map[string]interface{}) error {
	nstmt, err := d.PrepareNamed(query)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error select prepare named query %v", err)
		return err
	}

	err = nstmt.Select(list, conditionVal)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error select query %v", err)
		return err
	}
	return nil
}

// ValidateList -
func ValidateList(d *sqlx.DB, table string, column string, condition string, list interface{}) (bool, error) {
	values := map[string]interface{}{}
	queryValues := []string{}

	var total int64 = 0
	listInt64, ok := list.([]int64)
	if ok {
		for k, v := range listInt64 {
			col := "val_" + strconv.Itoa(k)
			values[col] = v
			queryValues = append(queryValues, ":"+col)
		}
		total = int64(len(listInt64))
	} else {
		listString, ok := list.([]string)
		if ok {
			for k, v := range listString {
				col := "val_" + strconv.Itoa(k)
				values[col] = v
				queryValues = append(queryValues, ":"+col)
			}
			total = int64(len(listString))
		} else {
			err := errors.New("Error validate list: interface type not supported")
			liblogger.Log(nil, true).Error(err)
			return false, err
		}
	}

	t := new(ModelTotal)
	_, err := Get(d, t, "SELECT COUNT("+column+") AS total FROM "+table+" WHERE "+column+" IN ("+strings.Join(queryValues, ", ")+") "+condition, values)
	if err != nil {
		return false, err
	}

	if t.Total == total {
		return true, err
	}
	return false, err
}

// Create - Create from query
func Create(d *sqlx.DB, cfg Config, dt interface{}, mode Mode, returnData bool) error {
	driver := d.DriverName()
	query, val := PrepareInsert(cfg.Table, dt, mode)
	if driver == "postgres" {
		if returnData {
			query += " RETURNING *"
		}
		_, err := Get(d, dt, query, val)
		return err
	}
	id, _, err := Exec(d, query, val)
	if returnData {
		_, err = GetByID(d, cfg, dt, id)
	}
	return err
}

// Update - General update from query
func Update(d *sqlx.DB, cfg Config, old interface{}, new interface{}, mode Mode, pk interface{}, returnData bool) (map[string]interface{}, error) {
	diff, err := UpdateCustom(d, cfg, old, new, mode, "AND id = :id ", map[string]interface{}{
		"id": pk,
	}, returnData)
	return diff, err
}

// UpdateCustom - Custome update from query
func UpdateCustom(d *sqlx.DB, cfg Config, old interface{}, new interface{}, mode Mode, condition string, conditionVal map[string]interface{}, returnData bool) (map[string]interface{}, error) {
	driver := d.DriverName()
	query, val, diff := PrepareUpdate(cfg.Table, old, new, condition, conditionVal, mode)
	if driver == "postgres" {
		if returnData {
			query += " RETURNING *"
		}
		_, err := Get(d, new, query, val)
		return diff, err
	}
	_, _, err := Exec(d, query, val)
	if returnData {
		_, err = GetByField(d, cfg, new, conditionVal, condition)
	}
	return diff, err
}

// Delete - General delete based on table configuration
func Delete(d *sqlx.DB, cfg Config, pk interface{}) error {
	if cfg.SoftDelete {
		return SoftDelete(d, cfg, pk)
	}
	return HardDelete(d, cfg, pk)
}

// HardDelete - General hard delete from query
func HardDelete(d *sqlx.DB, cfg Config, pk interface{}) error {
	return HardDeleteCustom(d, cfg, "AND id = :id ", map[string]interface{}{"id": pk})
}

// HardDeleteCustom - Custom hard delete from query
func HardDeleteCustom(d *sqlx.DB, cfg Config, condition string, conditionVal map[string]interface{}) error {
	query := "DELETE FROM " + cfg.Table + " WHERE 1=1 " + condition
	_, err := d.NamedExec(query, conditionVal)
	return err
}

// SoftDelete - General soft delete from query
func SoftDelete(d *sqlx.DB, cfg Config, pk interface{}) error {
	return SoftDeleteCustom(d, cfg, "AND id = :id ", map[string]interface{}{
		"id": pk,
	}, false)
}

// UnsoftDelete - General undo soft delete from query
func UnsoftDelete(d *sqlx.DB, cfg Config, pk interface{}) error {
	return SoftDeleteCustom(d, cfg, "AND id = :id ", map[string]interface{}{
		"id": pk,
	}, true)
}

// SoftDeleteCustom - Custom soft delete from query
func SoftDeleteCustom(d *sqlx.DB, cfg Config, condition string, conditionVal map[string]interface{}, revoke bool) error {
	delQuery := "deleted_at = "
	if revoke {
		delQuery += "NULL"
		condition += "AND deleted_at IS NOT NULL "
	} else {
		delQuery += ":deleted_at"
		condition += "AND deleted_at IS NULL "
		conditionVal["deleted_at"] = time.Now().UTC()
	}
	query := "UPDATE " + cfg.Table + " SET updated_at = NOW(), " + delQuery + " WHERE 1=1 " + condition
	_, err := d.NamedExec(query, conditionVal)
	return err
}

// GenerateUUID - Generate unique UUID
func GenerateUUID(d *sqlx.DB, cfg Config, dt interface{}) (string, error) {
	appMode := os.Getenv("app_mode")
	if appMode == "production" {
		appMode = ""
	} else {
		appMode = "dev-"
	}

	nUUID := uuid.NewString()
	exist, err := GetByUUID(d, cfg, dt, appMode+nUUID)
	if err != nil {
		return "", err
	}
	for exist {
		nUUID := uuid.NewString()
		exist, err = GetByUUID(d, cfg, dt, appMode+nUUID)
		if err != nil {
			return "", err
		}
	}
	return appMode + nUUID, nil
}
