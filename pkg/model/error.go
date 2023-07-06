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
	ErrCodeFileInvalid = CustomError{Code: 4005, Message: "ErrCodeFileInvalid",
		HttpCode: http.StatusBadRequest}
	ErrCodeRequestUserNotLogin = CustomError{Code: 4006, Message: "ErrCodeRequestUserNotLogin",
		HttpCode: http.StatusBadRequest}
	ErrCodePageSizeInvalid = CustomError{Code: 4007, Message: "ErrCodePageSizeInvalid",
		HttpCode: http.StatusBadRequest}
	ErrCodeInvalidUserRole = CustomError{Code: 4008, Message: "ErrCodeInvalidUserRole",
		HttpCode: http.StatusBadRequest}

	ErrCodeUnknownServerError = CustomError{Code: 5000, Message: "ErrCodeUnknownServerError",
		HttpCode: http.StatusInternalServerError}
	ErrCodeMysqlError = CustomError{Code: 5001, Message: "ErrCodeMysqlError",
		HttpCode: http.StatusInternalServerError}
	ErrCodeSessionError = CustomError{Code: 5002, Message: "ErrCodeSessionError",
		HttpCode: http.StatusInternalServerError}
	ErrCodeImageUploadError = CustomError{Code: 5003, Message: "ErrCodeImageUploadError",
		HttpCode: http.StatusInternalServerError}
	ErrCodeProfileNotFound = CustomError{Code: 5004, Message: "ErrCodeProfileNotFound",
		HttpCode: http.StatusInternalServerError}
	ErrCodeUnsupportEndorseOp = CustomError{Code: 5005, Message: "ErrCodeUnsupportEndorseOp",
		HttpCode: http.StatusInternalServerError}
)
