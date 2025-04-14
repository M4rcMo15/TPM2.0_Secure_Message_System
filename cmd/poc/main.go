package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Initialize reader for command line input
	reader := bufio.NewReader(os.Stdin)

	// Test for command line interface
	for {
		fmt.Print("Enter a command (type 'quit' to exit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "send":
			fmt.Print("Enter a message to send: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			fmt.Println("Message sent: ", message)

		case "receive":

			fmt.Println("Message received: ")
		case "quit":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Unknown command. Please try again.")
		}
	}
}
