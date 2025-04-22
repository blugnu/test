package json

import (
	"encoding/json"
	"reflect"

	"github.com/tidwall/sjson"
)

func Marshal(anything any) ([]byte, error) {
	anyV := reflect.ValueOf(anything)

	if anything == nil || isBytes(anyV) || isEmptySlice(anyV) {
		return json.Marshal(anything)
	}

	result := []byte{}
	err := walk(anything, func(v reflect.Value, marshaller bool, path string) error {
		bytes, err := parse(v, marshaller)
		if err != nil {
			return err
		}

		if result, err = sjson.SetRawBytes(result, path, bytes); err != nil {
			return err
		}

		return nil
	})

	return result, err
}
