package app

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

func Validate(i interface{}) (errors []string) {
	// Return errors if any
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	err := validate.Struct(i)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(
				errors,
				fmt.Sprintf("field: %s, value: %s, tag required: %s", e.Field(), e.Value(), e.Tag()),
			)
		}
	}
	return
}
