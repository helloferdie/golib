package libdb

import "database/sql"

// Config - Database table configuration
type Config struct {
	Table      string
	Fields     string
	SoftDelete bool
	Module     string
}

// GetConditionSoftDelete - Get condition for soft delete
func (cfg *Config) GetConditionSoftDelete() string {
	if cfg.SoftDelete {
		return "AND deleted_at IS NULL "
	}
	return ""
}

// ModelTimestamp - General struct for store timestamp
type ModelTimestamp struct {
	CreatedAt sql.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

// ModelID - General struct for query ID primary key
type ModelID struct {
	ID int64 `db:"id"`
}

// ModelTotal - General struct for query total items
type ModelTotal struct {
	Total int64 `db:"total"`
}

// ModelSummary - Generate struct to calculate summarize total data
type ModelSummary struct {
	Label string `db:"label"`
	Total int64  `db:"total"`
}

// ModelGetRequest - General request body for view/update/delete
type ModelGetRequest struct {
	ID int64 `json:"id" loc:"common." validate:"required"`
}

// ModelPaginationRequest - General request body for list (pagination)
type ModelPaginationRequest struct {
	ShowAll          bool
	Page             int64  `json:"page" loc:"common." validate:"required,numeric,min=1"`
	ItemsPerPage     int64  `json:"items_per_page" loc:"common." validate:"required,numeric,min=1,max=500"`
	OrderByField     string `json:"order_by_field" loc:"common."`
	OrderByDirection string `json:"order_by_direction" loc:"common."`
	OrderCustom      string
}
