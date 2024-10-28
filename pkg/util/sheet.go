package util

import (
	"reflect"
)

func ReflectGetFieldByTagOrName[K any](name string) *reflect.StructField {
	var k K
	typeOfK := reflect.TypeOf(k)
	for i := 0; i < typeOfK.NumField(); i++ {
		field := typeOfK.Field(i)
		if sheetName, ok := field.Tag.Lookup("sheets"); ok && sheetName == name {
			return &field
		}
	}
	if field, ok := typeOfK.FieldByName(name); ok {
		return &field
	}
	return nil
}
