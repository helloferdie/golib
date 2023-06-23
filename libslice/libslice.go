package libslice

import "reflect"

// GetTagSlice - Get slice tag string from struct interface
func GetTagSlice(t interface{}, tag string) []string {
	s := []string{}
	rType := reflect.TypeOf(t)
	rVal := reflect.ValueOf(t)

	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		if field.Anonymous {
			if field.Type.Kind() != reflect.Struct {
				continue
			}
			s = append(s, GetTagSlice(rVal.Field(i).Interface(), tag)...)
		} else {
			s = append(s, field.Tag.Get(tag))
		}
	}
	return s
}

// Contains - Check slice contains value of string
func Contains(a string, list []string) (int, bool) {
	for k, b := range list {
		if b == a {
			return k, true
		}
	}
	return -1, false
}

// ContainsInt64 - Check slice contains value of int64
func ContainsInt64(a int64, list []int64) (int, bool) {
	for k, b := range list {
		if b == a {
			return k, true
		}
	}
	return -1, false
}

// Unique - Return unique slice of string
func Unique(list []string) []string {
	keys := make(map[string]bool)
	tmpList := []string{}
	for _, v := range list {
		if _, ok := keys[v]; !ok {
			keys[v] = true
			tmpList = append(tmpList, v)
		}
	}
	return tmpList
}

// UniqueInt64 - Return unique slice of int64
func UniqueInt64(list []int64) []int64 {
	keys := make(map[int64]bool)
	tmpList := []int64{}
	for _, v := range list {
		if _, ok := keys[v]; !ok {
			keys[v] = true
			tmpList = append(tmpList, v)
		}
	}
	return tmpList
}
