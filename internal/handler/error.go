package handler

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Errors     map[string]string `json:"message"`
	StatusCode int               `json:"status_code"`
}

func Converter(err error) map[string]string {
	var validationErrors validator.ValidationErrors

	var result = make(map[string]string)

	if errors.As(err, &validationErrors) {
		for _, fe := range validationErrors {
			switch fe.Tag() {
			case "required":
				result[fe.Field()] = fmt.Sprintf("Поле %s обязательно", fe.Field())
			case "email":
				result[fe.Field()] = fmt.Sprintf("Поле %s должно быть валидным email", fe.Field())
			default:
				result[fe.Field()] = fmt.Sprintf("Ошибка в поле %s", fe.Field())
			}
		}
		return result
	}

	result["body"] = "пустое содержимое"
	return result
}
