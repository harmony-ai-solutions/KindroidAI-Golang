package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/harmony-ai-solutions/KindroidAI-Golang/client"
)

func main() {
	apiKey := os.Getenv("KINDROID_API_KEY")
	aiID := os.Getenv("KINDROID_AI_ID")

	if apiKey == "" {
		log.Fatal("KINDROID_API_KEY environment variable not set")
	}
	if aiID == "" {
		log.Fatal("KINDROID_AI_ID environment variable not set")
	}

	kindroidClient := client.NewKindroidAI(apiKey, aiID)

	// Setup User
	errUser := kindroidClient.SetupUserAndPermissions()
	if errUser != nil {
		log.Fatal("Failed to set up user")
	}
	fmt.Printf("Authenticated as UserID: %s\n", kindroidClient.UserID)

	ctx := context.Background()
	messages, err := kindroidClient.GetChatHistory(ctx, aiID, 10) // Get last 10 messages
	if err != nil {
		log.Fatalf("Failed to get chat history: %v", err)
	}

	fmt.Println("Chat History:")
	for _, msg := range messages {
		fmt.Printf("[%s] %s: %s\n", msg.GetTime().Format("2006-01-02 15:04:05"), msg.Sender, msg.Message)
	}
}
