// cmd/main.go
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"password-manager/internal/api"
	cipher_client "password-manager/internal/cipher"
	"password-manager/internal/client"
	"password-manager/internal/scheduler"
)

var buildVersion string
var buildDate string

func main() {

	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}

	fmt.Printf("password-keepre version: %s, buidl date: %s \n", buildVersion, buildDate)

	key, err := cipher_client.GenerateEncryptionKey()

	if err != nil {
		fmt.Errorf("unexpected error while generate key %s", err)
	}
	User := client.NewUser(key)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		s, err := scheduler.UserTokenServerSheduler(ctx, User)

		if s != "success" {
			err = scheduler.SaveUserKey(User)
			fmt.Println(err)
			return scheduler.UserTokenScheduler(ctx, User)
		}
		return err
	})

	fmt.Println("Interactive client started. Enter 'q' to exit. \n " +
		"Enter help to see the description of all commands")

	api.RunCommands(User)

}
