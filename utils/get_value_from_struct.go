package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func GetValueFromStruct(keyWithDots string, object interface{}) (interface{}, error) {
	keySlice := strings.Split(keyWithDots, ".")
	v := reflect.ValueOf(object)
	// iterate through field names ,ignore the first name as it might be the current instance name
	// you can make it recursive also if want to support types like slice,map etc along with struct
	for _, key := range keySlice[1:] {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if !v.IsValid() {
			return nil, nil
		}

		// we only accept structs
		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("only accepts structs; got %T", v)
		}

		v = v.FieldByName(key)
	}

	if !v.IsValid() {
		return nil, fmt.Errorf("invalid value found at: %s", keyWithDots)
	}
	return v.Interface(), nil
}
