package main

import (
	"password-manager/internal/app"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := app.Start(); err != nil {
		logrus.Fatalf("unsuccessful initilization app: %v", err)
	}
}
