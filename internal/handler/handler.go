package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"password-manager/internal/models"
	"password-manager/internal/service_errors"
)

type ServiceInterface interface {
	CreateUser(ctx context.Context, req models.Users) (string, error)
	GenerateUserToken(ctx context.Context, req models.Users) (string, error)
	SaveUserPassword(ctx context.Context, req models.Password) error
	GetUserPassword(ctx context.Context, name, UID string) (models.Password, error)
	UpdateUserSavedPassword(ctx context.Context, req models.NewPassword) error
	DeleteUserSavedPassword(ctx context.Context, name, UID string) error
	GetAllUserPasswords(ctx context.Context, UID string) ([]models.PasswordName, error)
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
		ctx := r.Context()

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		token, err := s.CreateUser(ctx, req)

		if err != nil {
			service_errors.HandleError(w, err)
			return
		}
		w.Header().Set("Authorization", token)
	}
}

func LoginUser(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Users
		ctx := r.Context()

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		token, err := s.GenerateUserToken(ctx, req)

		if err != nil {
			service_errors.HandleError(w, err)
			return
		}

		w.Header().Set("Authorization", token)
	}
}
