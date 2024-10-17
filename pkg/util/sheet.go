package util

import (
	"fmt"
	"reflect"

	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"google.golang.org/api/sheets/v4"
)

func ParseCandidateOffer[K any](valueRange *sheets.ValueRange) dto.SheetCandidateOffer {
	row := valueRange.Values[0]
	for colIdx, cellIf := range row {
		cell := cellIf.(string)
		fmt.Println("colIdx", colIdx)
		fmt.Println("cellIf", cellIf)
		fmt.Println("cell", cell)
	}

	return dto.SheetCandidateOffer{}
}

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
