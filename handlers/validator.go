package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type (
	ValidatorHandler struct {
		validator *validator.Validate
	}

	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}
)

var (
	VHandler ValidatorHandler
)

func NewValidator() ValidatorHandler {
	v := validator.New()
	VHandler = ValidatorHandler{validator: v}
	fmt.Println("[HANDLER] Validator Handler Initialized")
	return VHandler
}

func (h ValidatorHandler) Validate(data interface{}) []ErrorResponse {
	validationsError := []ErrorResponse{}

	if h.validator == nil {
		panic("[ERROR] Validator is not initialized")
	}

	err := h.validator.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var elem ErrorResponse
			elem.Error = true
			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			validationsError = append(validationsError, elem)
		}
	}

	return validationsError
}
