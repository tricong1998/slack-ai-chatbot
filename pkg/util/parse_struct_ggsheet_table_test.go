package util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseStructToSheetTable(t *testing.T) {
	type TestStruct struct {
		Name string `sheet:"Name"`
		Age  int    `sheet:"Age"`
	}

	data := []TestStruct{
		{Name: "John", Age: 30},
		{Name: "Jane", Age: 25},
		{Name: "Doe", Age: 35},
	}

	expectedHeader := []string{"Name", "Age"}
	expectedResult := [][]string{
		expectedHeader,
		{"John", "30"},
		{"Jane", "25"},
		{"Doe", "35"},
	}
	result := ParseStructToSheetTable(data)

	fmt.Println(result)

	if len(result) != len(data)+1 {
		t.Errorf("Expected %d rows, got %d", len(data)+1, len(result))
	}

	if !reflect.DeepEqual(result[0], expectedHeader) {
		t.Errorf("Expected header %v, got %v", expectedHeader, result[0])
	}

	for i, row := range result[0:] {
		if !reflect.DeepEqual(row, expectedResult[i]) {
			t.Errorf("Expected row %d %v, got %v", i, expectedResult[i], row)
		}
	}
}
