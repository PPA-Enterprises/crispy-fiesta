package event_log
import (
	"reflect"
	"fmt"
)

func structToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
		// we only accept structs
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("structToMap only accepts structs; got %T", val)
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		// gets us a StructField
		structField := typ.Field(i)
		if tagv := structField.Tag.Get(tag); tagv != "" {
			out[tagv] = val.Field(i).Interface()
		}
	}
	return out, nil
}

