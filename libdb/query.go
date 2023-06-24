package libdb

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/helloferdie/golib/libslice"
)

var defaultSkipColumn = []string{"id", "created_at", "updated_at", "deleted_at"}

// Mode - Database operation mode
//   - Skip: Set operation mode to skip columns in slice
//   - Only: Set operation mode to only applicable to columns in slice
//   - AutoTimestamp: Set `true` to let timestamp auto fill by database server, `false` to manual fill by code
type Mode struct {
	Skip          []string
	Only          []string
	AutoTimestamp bool
}

// DefaultMode -
var DefaultMode = Mode{}

// ModeAutoTimestamp -
var ModeAutoTimestamp = Mode{
	AutoTimestamp: true,
}

// PrepareInsert - Prepare insert query
func PrepareInsert(table string, data interface{}, mode Mode) (string, map[string]interface{}) {
	// Load mode
	m := "skip"
	var col, val, checkColumn []string
	if len(mode.Only) > 0 {
		m = "only"
		checkColumn = mode.Only
	} else {
		if len(mode.Skip) == 0 {
			checkColumn = defaultSkipColumn
		} else {
			checkColumn = mode.Skip
		}
	}

	// Map columns and value named parameters
	dataMap := MapTagDB(data, map[string]interface{}{})
	for tag := range dataMap {
		_, exist := libslice.Contains(tag, checkColumn)
		if (m == "only" && !exist) || (m == "skip" && exist) {
			delete(dataMap, tag)
			continue
		}

		col = append(col, "`"+tag+"`")
		val = append(val, ":"+tag)
	}

	// Manual assign timestamp
	if m == "skip" && len(mode.Skip) == 0 && !mode.AutoTimestamp {
		col = append(col, "created_at", "updated_at")
		val = append(val, ":created_at", ":updated_at")
		dataMap["created_at"] = time.Now().UTC()
		dataMap["updated_at"] = time.Now().UTC()
	}
	return "INSERT INTO " + table + " (" + strings.Join(col, ", ") + ") VALUES (" + strings.Join(val, ", ") + ")", dataMap
}

// PrepareUpdate - Prepare update query
func PrepareUpdate(table string, old interface{}, new interface{}, condition string, conditionVal map[string]interface{}, mode Mode) (string, map[string]interface{}, map[string]interface{}) {
	// Load mode
	m := "skip"
	var col, val, checkColumn []string
	if len(mode.Only) > 0 {
		m = "only"
		checkColumn = mode.Only
	} else {
		if len(mode.Skip) == 0 {
			checkColumn = defaultSkipColumn
		} else {
			checkColumn = mode.Skip
		}
	}

	// Map columns and value named parameters
	oldMap := MapTagDB(old, map[string]interface{}{})
	newMap := MapTagDB(new, map[string]interface{}{})
	dataMap := map[string]interface{}{}
	diffMap := map[string]interface{}{}
	for tag := range oldMap {
		_, exist := libslice.Contains(tag, checkColumn)
		if (m == "only" && !exist) || (m == "skip" && exist) {
			continue
		}

		oldVal := oldMap[tag]
		newVal := newMap[tag]
		if oldVal == newVal {
			continue
		}

		col = append(col, "`"+tag+"` = :"+tag)
		dataMap[tag] = newVal
		diffMap[tag] = map[string]interface{}{
			"o": oldVal,
			"n": newVal,
		}
	}

	// Manual assign timestamp
	if m == "skip" && len(mode.Skip) == 0 && !mode.AutoTimestamp {
		col = append(col, "updated_at = :updated_at")
		val = append(val, ":updated_at")
		dataMap["updated_at"] = time.Now().UTC()
	}

	for ck, cv := range conditionVal {
		dataMap[ck] = cv
	}
	return "UPDATE " + table + " SET " + strings.Join(col, ", ") + " WHERE 1=1 " + condition, dataMap, diffMap
}

// PrepareInQuery - Prepare query for in condition
func PrepareInQuery[T int64 | string](condition string, named string, list []T, values map[string]interface{}) string {
	query := []string{}
	for k, v := range list {
		s := named + "_" + strconv.Itoa(k)
		query = append(query, ":"+s)
		values[s] = v
	}

	if len(query) > 0 {
		return strings.TrimSpace(condition) + " " + "(" + strings.Join(query, ", ") + ")"
	}
	return ""
}

// PrepareOrderQuery - Prepare order query from pagination request
func PrepareOrderQuery(m *ModelPaginationRequest) (string, map[string]interface{}) {
	s := ""
	v := map[string]interface{}{}
	if m.OrderCustom != "" {
		s = m.OrderCustom
	} else if m.OrderByField != "" {
		s = "ORDER BY " + m.OrderByField + " " + m.OrderByDirection
	}
	s = strings.TrimSpace(s) + " "
	if !m.ShowAll {
		s += "LIMIT :limit OFFSET :offset"
		v["limit"] = m.ItemsPerPage
		v["offset"] = (m.Page - 1) * m.ItemsPerPage
	}
	return strings.TrimSpace(s), v
}

// MapTagDB -
func MapTagDB(t interface{}, result map[string]interface{}) map[string]interface{} {
	rVal := reflect.ValueOf(t)
	if rVal.Kind() == reflect.Ptr {
		rVal = rVal.Elem()
	}
	rType := rVal.Type()

	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		tag := field.Tag.Get("db")
		if tag != "" {
			result[tag] = rVal.Field(i).Interface()
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			result = MapTagDB(rVal.Field(i).Interface(), result)
		}
	}
	return result
}

// ModelCondition -
type ModelCondition struct {
	Query string
	Field []string
	Value []interface{}
}

// conditionBase -
func (mc *ModelCondition) conditionBase(column string, named string, t interface{}, mode string, condition string) {
	rVal := reflect.ValueOf(t)
	if !rVal.IsZero() {
		if named == "" {
			named = column
		}

		if condition == "LIKE" || condition == "NOT LIKE" {
			// Support postgres and mysql by tolower string
			tmp, _ := t.(string)
			tmp = "%" + strings.ToLower(tmp) + "%"
			mc.Value = append(mc.Value, tmp)
			column = "LOWER(" + column + ")"
		} else {
			mc.Value = append(mc.Value, t)
		}
		mc.Field = append(mc.Field, named)
		mc.Query += mode + " " + column + " " + condition + " :" + named + " "
	}
}

// Equal -
func (mc *ModelCondition) Equal(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "AND", "=")
}

// EqualCaseInsesitive -
func (mc *ModelCondition) EqualCaseInsesitive(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "AND", "LIKE")
}

// NotEqual -
func (mc *ModelCondition) NotEqual(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "AND", "!=")
}

// Like -
func (mc *ModelCondition) Like(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "AND", "LIKE")
}

// NotLike -
func (mc *ModelCondition) NotLike(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "AND", "NOT LIKE")
}

// OrEqual -
func (mc *ModelCondition) OrEqual(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "OR", "=")
}

// OrNotEqual -
func (mc *ModelCondition) OrNotEqual(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "OR", "!=")
}

// OrLike -
func (mc *ModelCondition) OrLike(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "OR", "LIKE")
}

// OrNotLike -
func (mc *ModelCondition) OrNotLike(column string, named string, t interface{}) {
	mc.conditionBase(column, named, t, "OR", "NOT LIKE")
}

// GetValue -
func (mc *ModelCondition) GetValue() map[string]interface{} {
	m := make(map[string]interface{}, len(mc.Field))
	for k, v := range mc.Field {
		m[v] = mc.Value[k]
	}
	return m
}

// ValueLike -
func ValueLike(val string) string {
	return "%" + strings.ToLower(val) + "%"
}
