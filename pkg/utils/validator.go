package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// CoreValidator is a wrapper around validator.Validate that adds the ability to register custom validators
type CoreValidator struct {
	validate *validator.Validate
}

// NewCoreValidator creates a new instance of CoreValidator
func NewCoreValidator() (*CoreValidator, error) {
	cv := &CoreValidator{
		validate: validator.New(),
	}
	if err := cv.RegisterValidators(map[string]func(level validator.FieldLevel) bool{
		"isYaml": cv.isYaml,
	}); err != nil {
		return nil, err
	}
	return cv, nil
}

// RegisterValidator registers a custom validator with the CoreValidator
func (v *CoreValidator) RegisterValidator(name string, validator func(level validator.FieldLevel) bool) error {
	return v.validate.RegisterValidation(name, validator)
}

// RegisterValidators registers a map of custom validators with the CoreValidator
func (v *CoreValidator) RegisterValidators(validators map[string]func(level validator.FieldLevel) bool) error {
	// Register the validators and return the first error if any
	for name, validatorFn := range validators {
		if err := v.RegisterValidator(name, validatorFn); err != nil {
			return err
		}
	}
	return nil
}

// Validate validates the given struct using the registered validators
func (v *CoreValidator) Validate(out interface{}) error {
	rv := reflect.ValueOf(out)

	// check to see if out is of kind struct
	if rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expecting struct, got %T instead", out)
	}

	if err := v.validate.Struct(out); err != nil {
		return v.validationError(err)
	}

	return nil
}

// Validator returns the underlying validator.Validate instance
func (v *CoreValidator) Validator() *validator.Validate {
	return v.validate
}

// isYaml validates the file extension of the given file name
func (v *CoreValidator) isYaml(fl validator.FieldLevel) bool {
	ext := filepath.Ext(fl.Field().String())
	return ext == ".yaml" || ext == ".yml"
}

// validationError converts a validator error into a more friendly error message
func (v *CoreValidator) validationError(err error) error {
	var ve validator.ValidationErrors

	// cast to validation errors and wrap the original error with a new one
	if errors.As(err, &ve) {
		for _, fe := range ve {
			return fmt.Errorf("%w: failed constraint %s=%s, received: %+v", err, fe.Tag(), fe.Param(), fe.Value())
		}
	}

	// if unable to cast err to a validator.ValidationErrors, then return the original error
	return err
}
