package utils

import "reflect"

func IsNil(v any) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).IsNil()
}
