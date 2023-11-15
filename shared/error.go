package shared

import (
	"fmt"
	"net/http"
)

var (
	ErrGettingJobs        = NewCustomError(http.StatusInternalServerError, "error getting all jobs")
	ErrCreatingJobs       = NewCustomError(http.StatusInternalServerError, "error creating jobs")
	ErrInvalidRequestBody = NewCustomError(http.StatusBadRequest, "invalid request body")
	ErrFindingJobs        = NewCustomError(http.StatusInternalServerError, "error finding jobs")
	ErrIdNotFound         = NewCustomError(http.StatusBadRequest, "id not found")
	ErrRecordNotFound     = NewCustomError(http.StatusBadRequest, "record not found")
	ErrCreateUsers        = NewCustomError(http.StatusInternalServerError, "error creating users")
	ErrInvalidToken       = NewCustomError(http.StatusUnauthorized, "error invalid token")
	ErrInvalidAuthHeader  = NewCustomError(http.StatusUnauthorized, "error invalid auth header")
	ErrJobNotFound        = NewCustomError(http.StatusBadRequest, "error job not found")
	ErrUnauthorized       = NewCustomError(http.StatusBadRequest, "error unauthorized")
	ErrMinusQuota         = NewCustomError(http.StatusBadRequest, "quota is less than zero")
	ErrJobTransaction     = NewCustomError(http.StatusInternalServerError, "error job transaction")
	ErrCreateApplyJob     = NewCustomError(http.StatusInternalServerError, "error creating apply job")
	ErrGettingUserJob     = NewCustomError(http.StatusInternalServerError, "error getting user job")
	ErrAlreadyApplied     = NewCustomError(http.StatusBadRequest, "already applied to the job")
	ErrUserDoesntExist    = NewCustomError(http.StatusBadRequest, "invalid email or password")
	ErrFailedLogin        = NewCustomError(http.StatusInternalServerError, "error failed login")
	ErrInvalidPassword    = NewCustomError(http.StatusBadRequest, "invalid email or password")
)

type CustomError struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

type ErrorDTO struct {
	Message string `json:"message"`
}

func NewCustomError(statuscode int, message string) *CustomError {
	return &CustomError{
		StatusCode: statuscode,
		Message:    message,
	}
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", ce.StatusCode, ce.Message)
}

func (ce *CustomError) ToErrorDTO() *ErrorDTO {
	return &ErrorDTO{
		Message: ce.Message,
	}
}
