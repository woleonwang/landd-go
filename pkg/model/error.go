package model

import "net/http"

type CustomError struct {
	Code     int
	Message  string
	HttpCode int
}

var (
	ErrCodeInvalidRequest = CustomError{Code: 4001, Message: "ErrCodeInvalidRequest",
		HttpCode: http.StatusBadRequest}
	ErrCodeDuplicateEmail = CustomError{Code: 4002, Message: "ErrCodeDuplicateEmail",
		HttpCode: http.StatusBadRequest}
	ErrCodeEmailNotFound = CustomError{Code: 4003, Message: "ErrCodeEmailNotFound",
		HttpCode: http.StatusBadRequest}
	ErrCodeIncorrectCredential = CustomError{Code: 4004, Message: "ErrCodeIncorrectCredential",
		HttpCode: http.StatusBadRequest}

	ErrCodeUnknownServerError = CustomError{Code: 5000, Message: "ErrCodeUnknownServerError",
		HttpCode: http.StatusInternalServerError}
	ErrCodeMysqlError = CustomError{Code: 5001, Message: "ErrCodeMysqlError",
		HttpCode: http.StatusInternalServerError}
	ErrCodeSessionError = CustomError{Code: 5002, Message: "ErrCodeSessionError",
		HttpCode: http.StatusInternalServerError}
)
