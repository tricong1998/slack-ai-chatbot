package util

import (
	"reflect"

	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
)

func GetStructFields(data interface{}) []string {
	fields := []string{}
	t := reflect.TypeOf(data)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fields = append(fields, field.Name)
	}
	return fields
}

func FromCandidateToEmployee(data *dto.SheetCandidateOffer) *dto.NewEmployeesSkills {

	return &dto.NewEmployeesSkills{
		FullName:       data.FullName,
		Position:       data.Position,
		Division:       data.Division,
		Level:          data.Level,
		OnBoardingDate: data.OnboardDate,
	}
}
