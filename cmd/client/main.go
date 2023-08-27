// cmd/main.go
package main

import (
	"client-password/internal/api"
	cipher_client "client-password/internal/cipher"
	"client-password/internal/client"
	"client-password/internal/scheduler"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func main() {
	key, err := cipher_client.GenerateEncryptionKey()

	if err != nil {
		fmt.Errorf("unexpected error while generate key %s", err)
	}
	User := client.NewUser(key)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return scheduler.UserTokenScheduler(ctx, User)
	})

	fmt.Println("Interactive client started. Enter 'q' or 'quit' to exit.")

	api.RunCommands(User)

}
