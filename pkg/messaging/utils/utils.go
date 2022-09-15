package utils

import (
	"github.com/iancoleman/strcase"
	"reflect"
)

func GetMessageName(message any) string {
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}
