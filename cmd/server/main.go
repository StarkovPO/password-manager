package main

import "github.com/sirupsen/logrus"

func main() {

	err := Start()
	if err != nil {
		logrus.Fatalf("unsuccessful initilization app: %v", err)
	}
}
