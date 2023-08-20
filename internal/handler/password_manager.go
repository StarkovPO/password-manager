package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"password-manager/internal/models"
	"password-manager/internal/service_errors"
)

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

		ctx := r.Context()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		req.UserID = r.Header.Get("User-ID")
		err = s.SaveUserPassword(ctx, req)

		if err != nil {

			service_errors.HandleError(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetUserPassword(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		if params["name"] == "" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		pas, err := s.GetUserPassword(r.Context(), params["name"], r.Header.Get("User-ID"))

		if err != nil {
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
