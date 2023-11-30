package libaudittrail

import "database/sql"

// Model -
type Model struct {
	ID         string       `db:"id" json:"id"`
	Operation  string       `db:"operation" json:"operation"`
	ModuleName string       `db:"module_name" json:"module_name"`
	TableName  string       `db:"table_name" json:"table_name"`
	TableKey   string       `db:"table_key" json:"table_key"`
	Change     string       `db:"change" json:"change"`
	Remark     string       `db:"remark" json:"remark"`
	ServiceIP  string       `db:"service_ip" json:"service_ip"`
	TokenID    string       `db:"token_id" json:"token_id"`
	CreatedBy  int64        `db:"created_by" json:"created_by"`
	CreatedAt  sql.NullTime `db:"created_at" json:"created_at"`
}
