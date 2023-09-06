package service_errors

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrBadRequest          = NewAppError(nil, "Bad request", "request in not supported format")
	ErrLoginAlreadyExist   = NewAppError(nil, "Login already exist", "User with this login already exist in DB")
	ErrCreateUser          = NewAppError(nil, "Something went wrong", "Executing query to create user ")
	ErrInvalidAuthHeader   = NewAppError(nil, "Invalid auth header", "Token do not contain Bearer or body")
	ErrInvalidLoginOrPass  = NewAppError(nil, "Invalid login or password", "User tried to login via incorrect login or pass")
	ErrEmptyNameOrPassword = NewAppError(nil, "Empty name or password in JSON", "User didn't fill require fields")
	ErrWithDB              = NewAppError(nil, "Unexpected service erorr", "DB returned unhandled error")
	ErrPasswordNotFound    = NewAppError(nil, "Password with that name not found", "No rows returned")
	ErrNameAlreadyExist    = NewAppError(nil, "The current name already exist in your walt, choose another", "Duplicate password name")
)

/*
AppError struct for custom errors
*/
type AppError struct {
	Err    error  `json:"-"`
	Msg    string `json:"msg"`
	DevMsg string `json:"dev_msg"`
}

// NewAppError create new app error
func NewAppError(err error, msg, devMsg string) *AppError {
	logrus.Errorf("Error: %v", devMsg)
	return &AppError{
		err,
		msg,
		devMsg,
	}
}

func (e *AppError) Error() string {
	return e.Msg
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

// HandleError handle error from app
func HandleError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	switch {
	case errors.Is(err, ErrBadRequest):
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrLoginAlreadyExist):
		http.Error(w, ErrLoginAlreadyExist.Error(), http.StatusConflict)
	case errors.Is(err, ErrCreateUser):
		http.Error(w, ErrCreateUser.Error(), http.StatusInternalServerError)
	case errors.Is(err, ErrInvalidAuthHeader):
		http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrInvalidLoginOrPass):
		http.Error(w, ErrInvalidLoginOrPass.Error(), http.StatusUnauthorized)
	case errors.Is(err, ErrWithDB):
		http.Error(w, ErrWithDB.Error(), http.StatusUnprocessableEntity)
	case errors.Is(err, ErrEmptyNameOrPassword):
		http.Error(w, ErrEmptyNameOrPassword.Error(), http.StatusConflict)
	case errors.Is(err, ErrPasswordNotFound):
		http.Error(w, ErrPasswordNotFound.Error(), http.StatusNotFound)
	default:
		logrus.Errorf("Unhandled error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
