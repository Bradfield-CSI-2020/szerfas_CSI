package sandbox

import (
	"fmt"
	"reflect"
)

func DeepEqual(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == b
	}

	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)

	return deepEqualValues(v1, v2)
}

func deepEqualValues(v1, v2 reflect.Value) bool {
	if v1.Type() != v2.Type() {
		return false
	}

	switch v1.Kind() {
	case reflect.Ptr:
		return deepEqualValues(v1.Elem(), v2.Elem())

	case reflect.Slice:
		if v1.Len() != v2.Len() {
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepEqualValues(v1.Index(i), v2.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Struct:
		for i := 0; i < v1.NumField(); i++ {
			if !deepEqualValues(v1.Field(i), v2.Field(i)) {
				return false
			}
		}
		return true

	case reflect.Map:
		if v1.Len() != v2.Len() {
			return false
		}
		iter := v1.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()

			if !deepEqualValues(v, v2.MapIndex(k)) {
				return false
			}
		}
		return true

	case reflect.Int:
		return v1.Int() == v2.Int()

	case reflect.String:
		return v1.String() == v2.String()

	case reflect.Bool:
		return v1.Bool() == v2.Bool()

	default:
		panic(fmt.Sprintf("Unexpected kind %v", v1.Kind()))
	}
}
