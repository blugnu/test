package jsonx

import "reflect"

func Marshal(v any) (string, error) {
	return visit(v), nil
}

func visit(v any) string {
	val := reflect.ValueOf(v)
	if val.IsNil() {
		return "null"
	}

	typ := reflect.TypeOf(v)
	switch typ.Kind() {
	case reflect.Struct:
		s := "{\n"
		for i := 0; i < typ.NumField(); i++ {
			fieldName := typ.Field(i).Name
			fieldVal := val.Field(i)
			s += fieldName + ":" + visit(fieldVal.Interface()) + "\n"
			if fieldVal.CanInterface() {
				visit(fieldVal.Interface())
			}
		}
		return s + "}"

	case reflect.Array, reflect.Slice:
		s := "["
		for i := 0; i < val.Len(); i++ {
			s += visit(val.Index(i).Interface())
			if i < val.Len()-1 {
				s += ",\n"
			}
		}
		return s + "]"

	case reflect.Map:
		s := "{"
		for _, key := range val.MapKeys() {
			keyStr := key.String()
			valStr := val.MapIndex(key).Interface()
			s += keyStr + ":" + visit(valStr) + ",\n"
		}
		return s + "}"
	}
	return "??"
}
