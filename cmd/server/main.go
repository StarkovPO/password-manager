package main

import (
	"password-manager/internal/app"

	"github.com/sirupsen/logrus"
)

// @title Password-manager API
// @version 1.0.0
// @description API server for password-manager CLI

// @host localhost:8080
// @BasePath /api

// @SecurityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := app.Start(); err != nil {
		logrus.Fatalf("unsuccessful initilization app: %v", err)
	}
}
