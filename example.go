// Package main
/*
Copyright © 2023 Harmony AI Solutions & Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	"fmt"
	"github.com/harmony-ai-solutions/KindroidAI-Golang/client"
	"os"
)

func main() {
	// Initial params - Usage of env vars recommended
	//apiKey := ""
	//kindroidId := ""
	apiKey := os.Getenv("KINDROID_API_KEY")
	kindroidId := os.Getenv("KINDROID_AI_ID")

	if apiKey == "" || kindroidId == "" {
		fmt.Println("Please set KINDROID_API_KEY and KINDROID_AI_ID environment variables.")
		os.Exit(1)
	}

	// Create client
	kindroidClient := client.NewKindroidAI(apiKey, kindroidId)

	// Initialize Input Scanner
	scanner := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Println()
		fmt.Print("Enter your message (type '!reset' to reset chat): ")
		scanner.Scan()

		if errInput := scanner.Err(); errInput != nil {
			fmt.Println(fmt.Errorf("unable to scan user input. Error: %v", errInput))
			os.Exit(2)
		}
		userInput := scanner.Text()

		// Check if the user wants to reset the chat
		if userInput == "!reset" {
			fmt.Print("Enter the AI's greeting message to start a new chat: ")
			scanner.Scan()
			greeting := scanner.Text()

			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading greeting: %v\n", err)
				os.Exit(3)
			}

			// Call ChatBreak with the greeting message
			if err := kindroidClient.ChatBreak(greeting); err != nil {
				fmt.Printf("Error resetting chat: %v\n", err)
				os.Exit(4)
			}

			fmt.Println("Chat has been reset.")
			fmt.Printf("AI: %s\n", greeting)
			continue
		}

		// Send the message to the AI
		response, err := kindroidClient.SendMessage(userInput)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			os.Exit(5)
		}

		// Display the AI's response
		fmt.Printf("AI: %s\n", response)
	}
}
