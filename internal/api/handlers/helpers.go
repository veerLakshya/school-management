package handlers

import (
	"errors"
	"reflect"
	"school-management/pkg/utils"
	"strings"
)

func CheckBlankFields(value interface{}) error {
	val := reflect.ValueOf(value)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			return utils.ErrorHandler(errors.New("All fields are required"), "All fields are required")
		}
	}
	return nil
}

func GetFieldNames(model interface{}) []string {
	val := reflect.TypeOf(model)
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldTag := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldTag) // get and append json tag for this field
	}
	return fields
}
