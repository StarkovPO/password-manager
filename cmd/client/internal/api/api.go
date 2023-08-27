package api

import (
	"bufio"
	cipher_client "client-password/internal/cipher"
	"client-password/internal/client"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	PasswordEndpoint = `/api/password`
)

func RunCommands(user *client.User) {
	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Print("Enter a command: ")
		scanner.Scan()
		command := scanner.Text()

		if command == "q" {
			fmt.Println("Exiting the client...")
			os.Exit(0)
		}

		// Split the command into parts
		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "sign-up":
			if len(parts) != 3 {
				fmt.Println("Usage: sign-up <username> <pass>")
				fmt.Println("Remove all spaces from login and password. Try again")
				continue
			}

			login := parts[1]
			password := parts[2]
			token, err := user.SignUp(login, password)
			if err != nil {
				fmt.Printf("Error: %s \n", err)
				continue
			}
			user.Token = token
			continue

		case "login":
			if len(parts) != 3 {
				fmt.Println("Usage: login <username> <pass>")
				fmt.Println("Remove all spaces from login and password. Try again")
				continue
			}
			login := parts[1]
			password := parts[2]
			token, err := user.Login(login, password)
			if err != nil {
				fmt.Printf("Error: %s", err)
				continue
			}
			user.Token = token
			continue
		case "save-pass":
			if len(parts) != 3 {
				fmt.Println("Usage: save-pass <name_pass> <pass>")
				fmt.Println("Remove all spaces from name of password and password. Try again")
				continue
			}

			encryptedPass, err := cipher_client.Encrypt(parts[2], user.EncryptionKey)
			if err != nil {
				fmt.Errorf("unexpected error while cipher the password: %s \n", err)
			}

			req := client.UserPass{
				Name:     parts[1],
				Password: encryptedPass,
			}

			_, err = user.Request(http.MethodPost, PasswordEndpoint, req)
			if err == nil {
				fmt.Println("Your password saved success")
				continue
			}
			fmt.Printf("Error: %s \n", err)
			continue
		case "get-pass":
			if len(parts) != 2 {
				fmt.Println("Usage: get-pass <name_pass>")
				fmt.Println("Remove all spaces from name of password. Try again")
				continue
			}

			res, err := user.Request(http.MethodGet, PasswordEndpoint, parts[1])

			if err != nil {
				fmt.Printf("Error: %s", err)
				continue
			}

			decryptedPass, err := cipher_client.Decrypt(res.(string), user.EncryptionKey)
			if err == nil {
				fmt.Printf("Your password: %s\n", decryptedPass)
				continue
			}
			fmt.Printf("Error: %s \n", err)
			continue

		default:
			fmt.Println("Unknown command. Available commands: sign-up, login, save-pass, q")
		}
	}
}
