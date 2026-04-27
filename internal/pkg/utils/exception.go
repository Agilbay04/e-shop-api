package utils

import "net/http"

type IAppError interface {
	GetStatusCode() int
}

type CustomError struct {
	StatusCode int
	Message    string
	Errors     interface{}
}

func (e *CustomError) GetStatusCode() int {
	return e.StatusCode
}

func (e *CustomError) Error() string {
	return e.Message
}

func BadRequestException(message string, errs interface{}) *CustomError {
	return &CustomError{StatusCode: http.StatusBadRequest, Message: message, Errors: errs}
}

func UnauthorizedException(message string) *CustomError {
	return &CustomError{StatusCode: http.StatusUnauthorized, Message: message}
}

func NotFoundException(message string) *CustomError {
	return &CustomError{StatusCode: http.StatusNotFound, Message: message}
}

func UnprocessableEntityException(message string) *CustomError {
	return &CustomError{StatusCode: http.StatusUnprocessableEntity, Message: message}
}

func InternalServerErrorException(message string) *CustomError {
    return &CustomError{StatusCode: http.StatusInternalServerError, Message: message}
}

func ForbiddenException(message string) *CustomError {
    return &CustomError{StatusCode: http.StatusForbidden, Message: message}
}

func ToManyRequestException(message string) *CustomError {
	return &CustomError{StatusCode: http.StatusTooManyRequests, Message: message}
}
