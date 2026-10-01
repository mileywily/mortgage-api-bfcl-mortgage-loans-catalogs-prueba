package error_fif

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"time"
)

const (
	description = "fail validation '%s' for '%s'"
)

func NewFifErrorHandler() ErrorHandler {
	return FifErrorHandler(DefaultErrorHandler())
}

func FifErrorHandler(next ErrorHandler) ErrorHandler {
	return func(err error) (int, interface{}) {
		switch err.(type) {
		case *time.ParseError:
			timeErr, _ := err.(*time.ParseError)
			return CreateBadRequest(fmt.Sprintf("error_handler time value: %s", timeErr.Value))

		case validator.ValidationErrors:
			typedErr, _ := err.(validator.ValidationErrors)
			var fv []interface{}
			for _, e := range typedErr {
				fv = append(fv, FieldViolation{Field: e.Namespace(), Description: fmt.Sprintf(description, e.Tag(), e.Value())})
			}
			return CreateBadRequest("validation error", fv...)

		case *json.SyntaxError:
			typedErr, _ := err.(*json.SyntaxError)
			return CreateBadRequest(fmt.Sprintf("json syntax error_handler at byte offset %d", typedErr.Offset))

		case *json.UnmarshalTypeError:
			typedErr, _ := err.(*json.UnmarshalTypeError)
			description := ""
			k := typedErr.Type.Kind()
			switch k {
			case reflect.Map:
				description = fmt.Sprintf("Error field '%v' expected object not %v", typedErr.Field, typedErr.Value)
			case reflect.Slice, reflect.Array:
				description = fmt.Sprintf("Error field '%v' expected array not %v", typedErr.Field, typedErr.Value)
			case reflect.Struct:
				description = fmt.Sprintf("Error field '%v' expected object not %v", typedErr.Field, typedErr.Value)
			default:
				description = fmt.Sprintf("Error field '%v' expected %v not %v", typedErr.Field, typedErr.Type, typedErr.Value)
			}

			var fv []interface{}
			fv = append(fv, FieldViolation{Field: typedErr.Field, Description: description})
			return CreateBadRequest("invalid json", fv...)

		default:
			return next(err)
		}
	}
}
