package libaudittrail

import (
	"encoding/json"

	"strconv"
	"time"

	"github.com/helloferdie/golib/libdb"
	"github.com/jmoiron/sqlx"
)

// generate - Generate default struct
func generate(cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) *Model {
	v, ok := key.(string)
	if !ok {
		i, ok := key.(int64)
		if ok {
			v = strconv.FormatInt(i, 10)
		} else {
			v = ""
		}
	}

	m := new(Model)
	m.TableName = cfg.Table
	m.TableKey = v
	m.ModuleName = cfg.Module
	m.Remark = remark
	m.TokenID = tokenID
	m.CreatedBy = creatorID
	return m
}

// GenerateID -
func (m *Model) GenerateID() {
	id, err := sf.NextID()
	for err != nil {
		id, err = sf.NextID()
	}
	m.ID = strconv.FormatUint(id, 10)
}

// Log -
func (m *Model) Log(d *sqlx.DB) error {
	loadConfig()

	m.ServiceIP = serviceIP
	m.CreatedAt.Valid = true
	m.CreatedAt.Time = time.Now().UTC()
	m.GenerateID()

	err := libdb.Create(d, TConfig, m, dbMode, false)
	return err
}

// LogCreate - Create record to database
func LogCreate(d *sqlx.DB, cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) error {
	bt, _ := json.Marshal(dt)

	m := generate(cfg, key, creatorID, tokenID, remark)
	m.Operation = "create"
	m.Change = string(bt)
	err := m.Log(d)
	return err
}

// LogUpdate - Update record from database
func LogUpdate(d *sqlx.DB, cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) error {
	bt, _ := json.Marshal(dt)

	m := generate(cfg, key, creatorID, tokenID, remark)
	m.Operation = "update"
	m.Change = string(bt)
	err := m.Log(d)
	return err
}

// LogDelete - Permanently delete record from database
func LogDelete(d *sqlx.DB, cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) error {
	bt, _ := json.Marshal(dt)

	m := generate(cfg, key, creatorID, tokenID, remark)
	m.Operation = "delete"
	m.Change = string(bt)
	err := m.Log(d)
	return err
}

// LogSoftDelete - Soft delete record from database
func LogSoftDelete(d *sqlx.DB, cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := generate(cfg, key, creatorID, tokenID, remark)
	m.Operation = "softdelete"
	err := m.Log(d)
	return err
}

// LogUnsoftDelete - Revert soft delete record from database
func LogUnsoftDelete(d *sqlx.DB, cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := generate(cfg, key, creatorID, tokenID, remark)
	m.Operation = "unsoftdelete"
	err := m.Log(d)
	return err
}

// LogView - View record from database
func LogView(d *sqlx.DB, cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := generate(cfg, key, creatorID, tokenID, remark)
	m.Operation = "view"
	err := m.Log(d)
	return err
}
