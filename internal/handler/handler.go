package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"password-manager/internal/models"
	"password-manager/internal/service_errors"
)

/*
ServiceInterface is interface for service module
*/
type ServiceInterface interface {
	CreateUser(ctx context.Context, req models.Users) (string, error)
	GenerateUserToken(ctx context.Context, req models.Users) (string, error)
	SaveUserPassword(ctx context.Context, req models.Password) error
	GetUserPassword(ctx context.Context, name, UID string) (models.Password, error)
	UpdateUserSavedPassword(ctx context.Context, req models.NewPassword) error
	DeleteUserSavedPassword(ctx context.Context, name, UID string) error
	GetAllUserPasswords(ctx context.Context, UID string) ([]models.PasswordName, error)
	SaveUserKey(ctx context.Context, UID string, key string) error
	GetUserKey(ctx context.Context, UID string) (string, error)
}

// @Summary Create User
// @tags Auth
// @Description register user in password manager
// @ID create-account
// @Accept json
// @Produce json
// @Param input body models.Users true "Creat the user with login and password"
// @Success 200
// @Failure 400,409 {object} service_errors.AppError
// @Failure 500 {object} service_errors.AppError
// @Failure default {object} service_errors.AppError
// @Router /api/user [post]
func RegisterUser(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Users

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		token, err := s.CreateUser(r.Context(), req)

		if err != nil {
			service_errors.HandleError(w, err)
			return
		}
		w.Header().Set("Authorization", token)
	}
}

// @Summary Login User
// @tags Auth
// @Description login user in password manager
// @ID login-account
// @Accept json
// @Produce json
// @Param input body models.Users true "Login the user with login and password"
// @Success 200
// @Failure 400,409 {object} service_errors.AppError
// @Failure 500 {object} service_errors.AppError
// @Failure default {object} service_errors.AppError
// @Router /api/login [post]
func LoginUser(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Users

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		token, err := s.GenerateUserToken(r.Context(), req)

		if err != nil {
			service_errors.HandleError(w, err)
			return
		}

		w.Header().Set("Authorization", token)
	}
}

func GetUserKey(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, service_errors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		key, err := s.GetUserKey(r.Context(), r.Header.Get("User-ID"))
		if err != nil {
			service_errors.HandleError(w, err)
			return
		}

		if key == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var UserKey models.UserKey
		UserKey.Key = key

		w.Header().Set("Content-type", "application/json")
		err = json.NewEncoder(w).Encode(&UserKey)

		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
	}
}

func SaveUserKey(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, service_errors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		body, _ := io.ReadAll(r.Body)
		parsedKey := string(body)

		err := s.SaveUserKey(r.Context(), r.Header.Get("User-ID"), parsedKey)
		if err != nil {
			service_errors.HandleError(w, err)
			return
		}
	}
}
