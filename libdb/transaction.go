package libdb

import (
	"time"

	"github.com/helloferdie/golib/liblogger"

	"github.com/jmoiron/sqlx"
)

// TxBegin - Begin database transaction connection
func TxBegin(d *sqlx.DB) (*sqlx.Tx, error) {
	tx, err := d.Beginx()
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error begin transaction connection %v", err)
	}
	return tx, err
}

// TxExec - Execute transaction query
func TxExec(tx *sqlx.Tx, query string, values map[string]interface{}) (int64, int64, error) {
	result, err := tx.NamedExec(query, values)
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

// TxGet - Get single row from transaction query
func TxGet(tx *sqlx.Tx, list interface{}, query string, values map[string]interface{}) (bool, error) {
	exist := false
	rows, err := tx.NamedQuery(query, values)
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

// TxGetByField - Get single row based on provided fields from transaction query
func TxGetByField(tx *sqlx.Tx, cfg Config, dt interface{}, params map[string]interface{}, condition string) (bool, error) {
	exist, err := TxGet(tx, dt, "SELECT "+cfg.Fields+" FROM "+cfg.Table+" WHERE 1=1 "+condition, params)
	return exist, err
}

// TxGetByID - Get single row by ID from transaction query
func TxGetByID(tx *sqlx.Tx, cfg Config, dt interface{}, id int64) (bool, error) {
	exist, err := TxGetByField(tx, cfg, dt, map[string]interface{}{
		"id": id,
	}, "AND id = :id "+cfg.GetConditionSoftDelete())
	return exist, err
}

// TxGetByUUID - Get single row by UUID from transaction query
func TxGetByUUID(tx *sqlx.Tx, cfg Config, dt interface{}, uuid string) (bool, error) {
	exist, err := TxGetByField(tx, cfg, dt, map[string]interface{}{
		"uuid": uuid,
	}, "AND uuid = :uuid "+cfg.GetConditionSoftDelete())
	return exist, err
}

// TxSelect - Select rows from transaction query
func TxSelect(tx *sqlx.Tx, list interface{}, query string, values map[string]interface{}) error {
	nstmt, err := tx.PrepareNamed(query)
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

// TxCreate - Create from transaction query
func TxCreate(tx *sqlx.Tx, cfg Config, dt interface{}, mode Mode, returnData bool) error {
	driver := tx.DriverName()
	query, val := PrepareInsert(cfg.Table, dt, mode)
	if driver == "postgres" {
		if returnData {
			query += " RETURNING *"
		}
		_, err := TxGet(tx, dt, query, val)
		return err
	}
	id, _, err := TxExec(tx, query, val)
	if returnData {
		_, err = TxGetByID(tx, cfg, dt, id)
	}
	return err
}

// TxUpdate - General update from transaction query
func TxUpdate(tx *sqlx.Tx, cfg Config, old interface{}, new interface{}, mode Mode, pk interface{}, returnData bool) (map[string]interface{}, error) {
	diff, err := TxUpdateCustom(tx, cfg, old, new, mode, "AND id = :id ", map[string]interface{}{
		"id": pk,
	}, returnData)
	return diff, err
}

// TxUpdateCustom - Custom update from transaction query
func TxUpdateCustom(tx *sqlx.Tx, cfg Config, old interface{}, new interface{}, mode Mode, condition string, conditionVal map[string]interface{}, returnData bool) (map[string]interface{}, error) {
	driver := tx.DriverName()
	query, val, diff := PrepareUpdate(cfg.Table, old, new, condition, conditionVal, mode)
	if driver == "postgres" {
		if returnData {
			query += " RETURNING *"
		}
		_, err := TxGet(tx, new, query, val)
		return diff, err
	}
	id, _, err := TxExec(tx, query, val)
	if returnData {
		_, err = TxGetByID(tx, cfg, new, id)
	}
	return diff, err
}

// TxDelete - General delete from transaction query
func TxDelete(tx *sqlx.Tx, cfg Config, pk interface{}) error {
	return TxDeleteCustom(tx, cfg, "AND id = :id ", map[string]interface{}{"id": pk})
}

// TxDeleteCustom - Custom delete from transaction query
func TxDeleteCustom(tx *sqlx.Tx, cfg Config, condition string, conditionVal map[string]interface{}) error {
	query := "DELETE FROM " + cfg.Table + " WHERE 1=1 " + condition
	_, err := tx.NamedExec(query, conditionVal)
	return err
}

// TxSoftDelete - General soft delete from transaction query
func TxSoftDelete(tx *sqlx.Tx, cfg Config, pk interface{}) error {
	return TxSoftDeleteCustom(tx, cfg, "AND id = :id ", map[string]interface{}{
		"id": pk,
	}, false)
}

// TxUnsoftDelete - General undo soft delete from transaction query
func TxUnsoftDelete(tx *sqlx.Tx, cfg Config, pk interface{}) error {
	return TxSoftDeleteCustom(tx, cfg, "AND id = :id ", map[string]interface{}{
		"id": pk,
	}, true)
}

// TxSoftDeleteCustom - Custom soft delete from transaction query
func TxSoftDeleteCustom(tx *sqlx.Tx, cfg Config, condition string, conditionVal map[string]interface{}, revoke bool) error {
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
	_, err := tx.NamedExec(query, conditionVal)
	return err
}
