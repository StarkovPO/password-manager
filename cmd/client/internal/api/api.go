package api

import (
	"bufio"
	"client-password/internal/client"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	savePasswordEndpoint = `/api/password`
)

func RunCommands(user client.User) {
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
				fmt.Printf("Error: %s", err)
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
			req := client.UserPass{
				Name:     parts[1],
				Password: parts[2],
			}

			_, err := user.Request(http.MethodPost, savePasswordEndpoint, req)
			if err == nil {
				fmt.Print("Your password saved success")
				continue
			}
			fmt.Printf("Error: %s \n", err)
			continue

		default:
			fmt.Println("Unknown command. Available commands: sign-up, login, save-pass, q")
		}
	}
}
