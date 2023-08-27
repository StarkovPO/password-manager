package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"password-manager/internal/models"
	"password-manager/internal/service_errors"
)

// @Summary Save Password
// @tags Password
// @Description save your password
// @Security ApiKeyAuth
// @ID save-pass
// @Accept json
// @Produce json
// @Param body body models.Password true "Creat the user with login and password"
// @Success 201
// @Failure 400,409 {object} service_errors.AppError
// @Failure 500 {object} service_errors.AppError
// @Failure default {object} service_errors.AppError
// @Router /api/password [post]
func SaveUserPassword(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Password

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, service_errors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		req.UserID = r.Header.Get("User-ID")
		err = s.SaveUserPassword(r.Context(), req)

		if err != nil {

			service_errors.HandleError(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary Get Password
// @tags Password
// @Description get your saved password
// @Security ApiKeyAuth
// @ID get-pass
// @Accept json
// @Produce json
// @Param name path string true "Search password by your name"
// @Success 200 {object} models.Password
// @Failure 400,409 {object} service_errors.AppError
// @Failure 500 {object} service_errors.AppError
// @Failure default {object} service_errors.AppError
// @Router /api/password/{name} [get]
func GetUserPassword(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, service_errors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		params := mux.Vars(r)
		if params["name"] == "" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		pas, err := s.GetUserPassword(r.Context(), params["name"], r.Header.Get("User-ID"))

		if err != nil {
			w.Header().Set("Content-type", "application/json")
			service_errors.HandleError(w, err)
			return
		}
		w.Header().Set("Content-type", "application/json")
		b, err := json.Marshal(pas)

		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			return
		}
	}
}

// @Summary Change saved Password
// @tags Password
// @Description change your password by name
// @Security ApiKeyAuth
// @ID change-pass
// @Accept json
// @Produce json
// @Param body body models.NewPassword true "Change the saved password"
// @Success 202
// @Failure 400,409 {object} service_errors.AppError
// @Failure 500 {object} service_errors.AppError
// @Failure default {object} service_errors.AppError
// @Router /api/password [put]
func UpdateUserSavedPassword(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.NewPassword

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, service_errors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		req.UserID = r.Header.Get("User-ID")

		err = s.UpdateUserSavedPassword(r.Context(), req)

		if err != nil {

			service_errors.HandleError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

	}
}

// @Summary Delete Password
// @tags Password
// @Description delete your saved password
// @Security ApiKeyAuth
// @ID del-pass
// @Accept json
// @Produce json
// @Param name path string true "Delete password by your name"
// @Success 204
// @Failure 400,409 {object} service_errors.AppError
// @Failure 500 {object} service_errors.AppError
// @Failure default {object} service_errors.AppError
// @Router /api/password/{name} [delete]
func DeleteUserSavedPassword(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, service_errors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		params := mux.Vars(r)
		if params["name"] == "" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		err := s.DeleteUserSavedPassword(r.Context(), params["name"], r.Header.Get("User-ID"))
		if err != nil {
			service_errors.HandleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary Get all names of your passwords
// @tags Password
// @Description get all nname of your saved password
// @Security ApiKeyAuth
// @ID get-pass-all
// @Accept json
// @Produce json
// @Success 200 {array} models.PasswordName
// @Failure 400,409 {object} service_errors.AppError
// @Failure 500 {object} service_errors.AppError
// @Failure default {object} service_errors.AppError
// @Router /api/password/all [get]
func GetAllUserPasswords(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, service_errors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		pass, err := s.GetAllUserPasswords(r.Context(), r.Header.Get("User-ID"))

		if err != nil {
			w.Header().Set("Content-type", "application/json")
			service_errors.HandleError(w, err)
			return
		}

		w.Header().Set("Content-type", "application/json")
		err = json.NewEncoder(w).Encode(&pass)

		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

	}
}
