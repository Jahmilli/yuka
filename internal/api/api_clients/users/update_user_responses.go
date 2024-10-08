// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"yuka/internal/api/api_models"
)

// UpdateUserReader is a Reader for the UpdateUser structure.
type UpdateUserReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateUserReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateUserOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateUserBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateUserInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[PUT /api/v1/users/{id}] updateUser", response, response.Code())
	}
}

// NewUpdateUserOK creates a UpdateUserOK with default headers values
func NewUpdateUserOK() *UpdateUserOK {
	return &UpdateUserOK{}
}

/*
UpdateUserOK describes a response with status code 200, with default header values.

OK
*/
type UpdateUserOK struct {
	Payload *api_models.ModelsUser
}

// IsSuccess returns true when this update user o k response has a 2xx status code
func (o *UpdateUserOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update user o k response has a 3xx status code
func (o *UpdateUserOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update user o k response has a 4xx status code
func (o *UpdateUserOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update user o k response has a 5xx status code
func (o *UpdateUserOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update user o k response a status code equal to that given
func (o *UpdateUserOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the update user o k response
func (o *UpdateUserOK) Code() int {
	return 200
}

func (o *UpdateUserOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/users/{id}][%d] updateUserOK %s", 200, payload)
}

func (o *UpdateUserOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/users/{id}][%d] updateUserOK %s", 200, payload)
}

func (o *UpdateUserOK) GetPayload() *api_models.ModelsUser {
	return o.Payload
}

func (o *UpdateUserOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(api_models.ModelsUser)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateUserBadRequest creates a UpdateUserBadRequest with default headers values
func NewUpdateUserBadRequest() *UpdateUserBadRequest {
	return &UpdateUserBadRequest{}
}

/*
UpdateUserBadRequest describes a response with status code 400, with default header values.

Bad Request
*/
type UpdateUserBadRequest struct {
	Payload *api_models.ModelsValidationError
}

// IsSuccess returns true when this update user bad request response has a 2xx status code
func (o *UpdateUserBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update user bad request response has a 3xx status code
func (o *UpdateUserBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update user bad request response has a 4xx status code
func (o *UpdateUserBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this update user bad request response has a 5xx status code
func (o *UpdateUserBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this update user bad request response a status code equal to that given
func (o *UpdateUserBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the update user bad request response
func (o *UpdateUserBadRequest) Code() int {
	return 400
}

func (o *UpdateUserBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/users/{id}][%d] updateUserBadRequest %s", 400, payload)
}

func (o *UpdateUserBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/users/{id}][%d] updateUserBadRequest %s", 400, payload)
}

func (o *UpdateUserBadRequest) GetPayload() *api_models.ModelsValidationError {
	return o.Payload
}

func (o *UpdateUserBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(api_models.ModelsValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateUserInternalServerError creates a UpdateUserInternalServerError with default headers values
func NewUpdateUserInternalServerError() *UpdateUserInternalServerError {
	return &UpdateUserInternalServerError{}
}

/*
UpdateUserInternalServerError describes a response with status code 500, with default header values.

Internal Server Error
*/
type UpdateUserInternalServerError struct {
	Payload *api_models.ModelsBaseError
}

// IsSuccess returns true when this update user internal server error response has a 2xx status code
func (o *UpdateUserInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update user internal server error response has a 3xx status code
func (o *UpdateUserInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update user internal server error response has a 4xx status code
func (o *UpdateUserInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update user internal server error response has a 5xx status code
func (o *UpdateUserInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update user internal server error response a status code equal to that given
func (o *UpdateUserInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the update user internal server error response
func (o *UpdateUserInternalServerError) Code() int {
	return 500
}

func (o *UpdateUserInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/users/{id}][%d] updateUserInternalServerError %s", 500, payload)
}

func (o *UpdateUserInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/users/{id}][%d] updateUserInternalServerError %s", 500, payload)
}

func (o *UpdateUserInternalServerError) GetPayload() *api_models.ModelsBaseError {
	return o.Payload
}

func (o *UpdateUserInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(api_models.ModelsBaseError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
