package util

import (
	"fmt"
	"reflect"
	"strings"
)

func ParseStructToSheetTable[K any](data []K) [][]interface{} {
	if len(data) == 0 {
		return [][]interface{}{}
	}

	// Get the type of the struct
	t := reflect.TypeOf(data[0])

	// Create the header row
	header := make([]interface{}, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("sheet")
		if tag != "" {
			header = append(header, strings.Split(tag, ",")[0])
		}
	}

	// Initialize the result with the header
	result := [][]interface{}{header}

	// Add data rows
	for _, item := range data {
		row := make([]interface{}, len(header))
		v := reflect.ValueOf(item)
		for i := range header {
			field := v.Field(i)
			row[i] = fmt.Sprintf("%v", field.Interface())
		}
		result = append(result, row)
	}

	return result
}
