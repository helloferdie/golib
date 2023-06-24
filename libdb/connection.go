package libdb

import (
	"fmt"
	"os"

	"github.com/helloferdie/golib/liblogger"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Connection -
type Connection struct {
	Driver string
	DSN    string
}

// cacheConnection - Cache connection string in memory
var cacheConnection = map[string]*Connection{}

// setConnection - Set connection string
func setConnection(env string) (*Connection, error) {
	if env == "" {
		env = "db"
	}

	v, ok := cacheConnection[env]
	if !ok {
		driver := os.Getenv(env + "_driver")
		host := os.Getenv(env + "_host")
		port := os.Getenv(env + "_port")
		user := os.Getenv(env + "_user")
		pass := os.Getenv(env + "_pass")
		dbname := os.Getenv(env + "_name")

		if driver == "mysql" {
			cfg := mysql.NewConfig()
			cfg.Net = "tcp"
			cfg.Addr = host + ":" + port
			cfg.User = user
			cfg.Passwd = pass
			cfg.DBName = dbname
			cfg.ParseTime = true
			cfg.Params = map[string]string{
				"charset": "utf8mb4",
			}

			cacheConnection[env] = &Connection{
				Driver: driver,
				DSN:    cfg.FormatDSN(),
			}
			return cacheConnection[env], nil
		}
		return nil, fmt.Errorf("Database driver not supported for %s", env)
	}
	return v, nil
}

// Open - Open connection with default retry parameter
func Open(env string) (*sqlx.DB, error) {
	return OpenRetry(env, 3)
}

// OpenRetry - Open connection with custom retry parameter
func OpenRetry(env string, maxRetry int) (*sqlx.DB, error) {
	conn, err := setConnection(env)
	if err != nil {
		liblogger.Log(nil, true).Errorf("Error set connection string %v", err)
		return nil, err
	}

	if maxRetry < 0 {
		maxRetry = 0
	}

	db, err := sqlx.Connect(conn.Driver, conn.DSN)
	if err != nil {
		return nil, err
	}
	return db, err
}
