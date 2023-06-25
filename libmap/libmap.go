package libmap

import "reflect"

// MapTag - Return map struct based given tag in struct
func MapTag(t interface{}, tag string, result map[string]interface{}) map[string]interface{} {
	rVal := reflect.ValueOf(t)
	if rVal.Kind() == reflect.Ptr {
		rVal = rVal.Elem()
	}
	rType := rVal.Type()

	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		tag := field.Tag.Get(tag)
		if tag != "" {
			result[tag] = rVal.Field(i).Interface()
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			result = MapTag(rVal.Field(i).Interface(), tag, result)
		}
	}
	return result
}
