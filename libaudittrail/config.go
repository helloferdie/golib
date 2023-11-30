package libaudittrail

import (
	"os"
	"strings"
	"time"

	"github.com/helloferdie/golib/libdb"
	"github.com/helloferdie/golib/libslice"
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake
var initialize = false
var serviceIP = ""
var dbMode = libdb.Mode{Skip: []string{"updated_at", "deleted_at"}, AutoTimestamp: true}

func init() {
	loc, _ := time.LoadLocation("UTC")
	sfTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-01-01 00:00:00", loc)

	var st sonyflake.Settings
	st.StartTime = sfTime
	st.MachineID = func() (uint16, error) {
		return 1, nil
	}
	sf = sonyflake.NewSonyflake(st)
}

// loadConfig -
func loadConfig() {
	if !initialize {
		name, err := os.Hostname()
		if err == nil {
			serviceIP = name + ":" + os.Getenv("port")
		}
		initialize = true
	}
}

// TConfig -
var TConfig = libdb.Config{
	Table:      "audit_trail",
	Fields:     strings.Join(libslice.GetTagSlice(Model{}, "db"), ", "),
	SoftDelete: true,
}
