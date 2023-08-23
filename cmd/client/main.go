// cmd/main.go
package main

import (
	"client-password/internal/api"
	"client-password/internal/client"
	"fmt"
)

//var rootCmd = &cobra.Command{
//	Use:   "client",
//	Short: "Interactive CLI client",
//	Run: func(cmd *cobra.Command, args []string) {
//		fmt.Println("Interactive client started. Enter 'q' or 'quit' to exit.")
//		stopChan := make(chan bool)
//
//		go func() {
//			for {
//				fmt.Print("> ")
//				var input string
//				fmt.Scanln(&input)
//
//				if strings.ToLower(input) == "q" || strings.ToLower(input) == "quit" {
//					stopChan <- true
//					return
//				}
//
//				// Execute Cobra commands
//				cmd.SetArgs(strings.Fields(input))
//				if err := cmd.Execute(); err != nil {
//					fmt.Println("Error:", err)
//				}
//			}
//		}()
//
//		sigChan := make(chan os.Signal, 1)
//		signal.Notify(sigChan, syscall.SIGINT)
//
//		select {
//		case <-stopChan:
//			fmt.Println("Exiting...")
//		case <-sigChan:
//			fmt.Println("Received Ctrl+C signal, exiting...")
//		}
//
//		fmt.Println("Client has exited.")
//	},
//}

func main() {

	User := client.NewUser()

	fmt.Println("Interactive client started. Enter 'q' or 'quit' to exit.")

	api.RunCommands(*User)

}
