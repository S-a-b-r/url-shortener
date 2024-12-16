package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK         = 200
	StatusBadRequest = 400
)

func Ok(data interface{}) Response {
	return Response{Status: StatusOK}
}

func Error(msg string) Response {
	return Response{Status: StatusBadRequest, Error: msg}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is not valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.ActualTag()))
		}
	}

	return Response{Status: StatusBadRequest, Error: strings.Join(errMsgs, ", ")}
}
