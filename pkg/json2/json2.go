package json2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Map map[string]json.RawMessage

func (m Map) Decode(output interface{}) error {
	value := reflect.ValueOf(output).Elem()
	structType := value.Type()
	fieldCount := structType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := structType.Field(i)
		jsonKey := strings.Split(field.Tag.Get("json"), ",")
		if raw, ok := m[jsonKey[0]]; ok {
			if err := json.Unmarshal(raw, value.Field(i).Addr().Interface()); err != nil {
				return fmt.Errorf("cannot decode JSON field '%s':\n%w", field.Name, err)
			}
		}
	}
	return nil
}

func (m Map) Encode(input interface{}) error {
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	structType := value.Type()
	fieldCount := structType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := structType.Field(i)
		v := value.Field(i)
		jsonKey := strings.Split(field.Tag.Get("json"), ",")
		if len(jsonKey) >= 2 && jsonKey[1] == "omitempty" && v.IsZero() {
			delete(m, jsonKey[0])
		} else {
			msg, err := json.Marshal(v.Interface())
			if err != nil {
				return fmt.Errorf("cannot encode JSON field '%s':\n%w", field.Name, err)
			}
			m[jsonKey[0]] = json.RawMessage(msg)
		}
	}
	return nil
}
