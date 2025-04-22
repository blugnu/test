package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parse(v reflect.Value, marshaller bool) ([]byte, error) {
	if !marshaller {
		return defaultParser(v)
	}
	return marshalJSON(v)
}

func marshalJSON(v reflect.Value) ([]byte, error) {
	vPtr := reflect.NewAt(v.Elem().Type(), v.UnsafePointer())
	results := vPtr.Elem().MethodByName("MarshalJSON").Call([]reflect.Value{})

	bytes, _ := results[0].Interface().([]byte)
	err, _ := results[1].Interface().(error)

	return bytes, err
}

func defaultParser(v reflect.Value) ([]byte, error) {
	s := fmt.Sprint(reflect.Indirect(v))

	if isBytes(v) {
		return parseBytes(v, s)
	}

	result, err := resolve(v, s)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func resolve(v reflect.Value, s string) (any, error) {
	vKind := v.Kind()

	if matchKind(vKind,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64) {
		return strconv.Atoi(s)
	}

	if matchKind(vKind,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64) {
		return strconv.ParseUint(s, 10, 64)
	}

	if matchKind(vKind, reflect.Float32, reflect.Float64) {
		return strconv.ParseFloat(s, 64)
	}

	if matchKind(vKind, reflect.Complex64, reflect.Complex128) {
		return strconv.ParseComplex(s, 128)
	}

	return s, nil
}

func isBytes(v reflect.Value) bool {
	return fmt.Sprint(v.Type()) == `[]uint8`
}

func isEmptySlice(v reflect.Value) bool {
	return matchKind(v.Kind(), reflect.Slice, reflect.Array) && v.Len() == 0
}

func matchKind(kind reflect.Kind, kinds ...reflect.Kind) bool {
	for i := range kinds {
		if kind == kinds[i] {
			return true
		}
	}
	return false
}

func parseBytes(v reflect.Value, s string) ([]byte, error) {
	if v.IsZero() {
		return []byte(`null`), nil
	}

	bytes := make([]byte, 0)
	chunks := strings.Split(s[1:len(s)-1], " ")

	for i := range chunks {
		val, err := strconv.ParseUint(chunks[i], 10, 64)
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, uint8(val))
	}

	return json.Marshal(bytes)
}
