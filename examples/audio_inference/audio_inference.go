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
	messages, err := kindroidClient.GetChatHistory(ctx, aiID, 10) // Get last messages
	if err != nil {
		log.Fatalf("Failed to get chat history: %v", err)
	}

	// Fetch first AI message we find
	var kindroidMessageForAudio *client.ChatMessage
	for _, msg := range messages {
		if msg.Sender != "ai" {
			continue
		}
		fmt.Println("Found AI message for Inference Test")
		fmt.Printf("[%s] %s: %s\n", msg.GetTime().Format("2006-01-02 15:04:05"), msg.Sender, msg.Message)
		kindroidMessageForAudio = msg
		break
	}

	if kindroidMessageForAudio == nil {
		log.Fatal("No AI message found for Inference Test")
	}

	audioBytes, errGenAudio := kindroidClient.AudioInference(kindroidMessageForAudio.ID)
	if errGenAudio != nil {
		log.Fatalf("Failed to generate audio: %v", errGenAudio)
	}
	fmt.Println(fmt.Sprintf("Generated Audio for Inference Test. Audio file size (bytes): %d", len(audioBytes)))

}
