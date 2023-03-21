package utils

import "reflect"

func Ptr[T any](v T) *T {
	return &v
}

func Unptr[T any](v *T) T {
	if v == nil {
		return reflect.Zero(reflect.TypeOf(v).Elem()).Interface().(T)
	}
	return *v
}
