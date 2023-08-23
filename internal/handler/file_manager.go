package handler

import (
	"io"
	"net/http"
	"password-manager/internal/models"
	"password-manager/internal/service_errors"
)

func SaveUserFile(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.File

		if r.Header.Get("Content-Type") != "multipart/form-data" {
			http.Error(w, service_errors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, service_errors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing file data", http.StatusInternalServerError)
			return
		}

		file, _, err := r.FormFile("file")

		if err != nil {
			http.Error(w, "Error retrieving file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file data", http.StatusInternalServerError)
			return
		}

		req.FileData = fileData
		req.Name = r.FormValue("name")
		req.UserID = r.Header.Get("User-ID")

	}
}
