// cmd/main.go
package main

import (
	"client-password/internal/api"
	cipher_client "client-password/internal/cipher"
	"client-password/internal/client"
	"fmt"
)

func main() {
	key, err := cipher_client.GenerateEncryptionKey()

	if err != nil {
		fmt.Errorf("unexpected error while generate key %s", err)
	}
	User := client.NewUser(key)

	fmt.Println("Interactive client started. Enter 'q' or 'quit' to exit.")

	api.RunCommands(*User)

}
