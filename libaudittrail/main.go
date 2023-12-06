package libaudittrail

import (
	"encoding/json"

	"strconv"
	"time"

	"github.com/helloferdie/golib/libdb"
	"github.com/jmoiron/sqlx"
)

// generate - Generate default struct
func generate(cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) *Model {
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

	if dt != nil {
		bt, _ := json.Marshal(dt)
		m.Change = string(bt)
	}
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

// PrepareLogCreate - Prepare create log
func PrepareLogCreate(cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) *Model {
	m := generate(cfg, dt, key, creatorID, tokenID, remark)
	m.Operation = "create"
	return m
}

// LogCreate - Create record to database
func LogCreate(d *sqlx.DB, cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := PrepareLogCreate(cfg, dt, key, creatorID, tokenID, remark)
	err := m.Log(d)
	return err
}

// PrepareLogUpdate - Prepare update log
func PrepareLogUpdate(cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) *Model {
	m := generate(cfg, dt, key, creatorID, tokenID, remark)
	m.Operation = "update"
	return m
}

// LogUpdate - Update record from database
func LogUpdate(d *sqlx.DB, cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := PrepareLogUpdate(cfg, dt, key, creatorID, tokenID, remark)
	err := m.Log(d)
	return err
}

// PrepareLogDelete - Prepare delete log
func PrepareLogDelete(cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) *Model {
	m := generate(cfg, dt, key, creatorID, tokenID, remark)
	m.Operation = "delete"
	return m
}

// LogDelete - Permanently delete record from database
func LogDelete(d *sqlx.DB, cfg libdb.Config, dt interface{}, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := PrepareLogDelete(cfg, dt, key, creatorID, tokenID, remark)
	err := m.Log(d)
	return err
}

// PrepareLogSoftDelete - Prepare soft delete log
func PrepareLogSoftDelete(cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) *Model {
	m := generate(cfg, nil, key, creatorID, tokenID, remark)
	m.Operation = "softdelete"
	return m
}

// LogSoftDelete - Soft delete record from database
func LogSoftDelete(d *sqlx.DB, cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := PrepareLogSoftDelete(cfg, key, creatorID, tokenID, remark)
	err := m.Log(d)
	return err
}

// PrepareLogUnsoftDelete - Prepare unsoft delete log
func PrepareLogUnsoftDelete(cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) *Model {
	m := generate(cfg, nil, key, creatorID, tokenID, remark)
	m.Operation = "unsoftdelete"
	return m
}

// LogUnsoftDelete - Revert soft delete record from database
func LogUnsoftDelete(d *sqlx.DB, cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := PrepareLogUnsoftDelete(cfg, key, creatorID, tokenID, remark)
	err := m.Log(d)
	return err
}

// PrepareLogView - Prepare view log
func PrepareLogView(cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) *Model {
	m := generate(cfg, nil, key, creatorID, tokenID, remark)
	m.Operation = "view"
	return m
}

// LogView - View record from database
func LogView(d *sqlx.DB, cfg libdb.Config, key interface{}, creatorID int64, tokenID string, remark string) error {
	m := PrepareLogView(cfg, key, creatorID, tokenID, remark)
	err := m.Log(d)
	return err
}
