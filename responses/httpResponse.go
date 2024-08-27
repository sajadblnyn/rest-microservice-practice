package responses

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Field string
	Msg   string
}
type HttpResponse struct {
	Status string      `json:"status"`
	Errors []ApiError  `json:"errors"`
	Data   interface{} `json:"data"`
}

func MakeResponse(err error, data interface{}) *HttpResponse {
	status := "Success"
	if err != nil {
		status = "Failed"
	}

	var ve validator.ValidationErrors
	outErrors := make([]ApiError, len(ve))

	if errors.As(err, &ve) {
		for _, fe := range ve {
			outErrors = append(outErrors, ApiError{fe.Field(), msgForTag(fe.Tag())})
		}
	}

	return &HttpResponse{Status: status, Errors: outErrors}
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "titleLength":
		return "Invalid title"
	}
	return ""
}
