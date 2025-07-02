# 💬 KindroidAI - Golang Client by Project Harmony.AI
![Tag](https://img.shields.io/github/license/harmony-ai-solutions/CharacterAI-Golang)

An unofficial API Client for KindroidAI, written in Golang.

Based on KindroidAI's public API reference: [KindroidAI API](https://docs.kindroid.ai/api-documentation)

---

## ⚠️ ATTENTION: Pre-Release!
This repo is currently very barebone and not optimized for productive usage in golang applications yet.
Any support with testing, verifying functionality and adding proper golang struct handling is heavily appreciated.

You have questions, need help, or just want to show your support? Reach us here: [Discord Server & Patreon page](#how-to-reach-out-to-us).

### TODO's:
- [x] Implement API accodring to public documentation
    - [x] Confirm basic functionality
    - [x] Write tests for all API Methods

## 💻 Installation
```bash
go get github.com/harmony-ai-solutions/KindroidAI-Golang
```

## 📚 Documentation
Detailed documentation and Apidocs coming soon.

### Experimental Features (Use at Your Own Risk)
The following API methods and features have been discovered through network analysis (HAR logs) and are not part of the official KindroidAI API documentation. They may change or be removed without notice. Use them at your own risk in production environments.

- **`CheckUserSubscription()`**: Retrieves detailed user subscription information.
- **`AudioInference(messageID string)`**: Sends a request related to audio processing for a given message.
- **`GetChatHistory(ctx context.Context, aiID string, limit int) ([]ChatMessage, error)`**: Retrieves the chat history for a given AI. This method communicates with Google's Firestore like the Kindroid Client does and decrypts the messages.
- **Enhanced `SendMessage` options**: The `SendMessageAdvanced` method and its `SendMessageOptions` struct expose additional parameters (e.g., `ImageURLs`, `VideoURL`, `Stream`) that are not explicitly documented in the public API reference.

## 📙 Examples

### Authentication
The client can be authenticated in two ways:

1.  **Using a JWT (Bearer Token)**: You can provide the short-lived bearer token obtained from the web application's network traffic as the `KINDROID_API_KEY`. The client will automatically parse the token to extract your `UserID`, which is required for fetching chat history.
2.  **Using a Static API Key**: If you are using a permanent API key from your Kindroid account settings, you must also provide your `UserID` separately via the `KINDROID_USER_ID` environment variable. This is necessary because the static key does not contain the user ID.

### Basic Chat App
Example code for a simple, functional Chat app. The code can also be found in [example.go](example.go)
```Golang
package main

import (
  "bufio"
  "fmt"
  "github.com/harmony-ai-solutions/KindroidAI-Golang/client"
  "os"
)

func main() {
  // Initial params - Usage of env vars recommended
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
```

### Fetching Chat History
Example code for fetching and decrypting chat history. The code can also be found in [examples/example_chat_history.go](examples/example_chat_history.go)
```go
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

	// Ensure UserID is extracted (it's done in NewKindroidAI)
	// Ensure UserID is available for chat history features
	if kindroidClient.UserID == "" {
		log.Fatal("UserID not found. If using a static API key, ensure KINDROID_USER_ID is set. Otherwise, provide a valid JWT.")
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
```

---

## About Project Harmony.AI
![Project Harmony.AI](docs/images/Harmony-Main-Banner-200px.png)
### Our goal: Elevating Human <-to-> AI Interaction beyond known boundaries.
Project Harmony.AI emerged from the idea to allow for a seamless living together between AI-driven characters and humans.
Since it became obvious that a lot of technologies required for achieving this goal are not existing or still very experimental,
the long term vision of Project Harmony is to establish the full set of technologies which help minimizing biological and
technological barriers in Human <-to-> AI Interaction.

### Our principles: Fair use and accessibility
We want to counter today's tendencies of AI development centralization at the hands of big
corporations. We're pushing towards maximum transparency in our own development efforts, and aim for our software to be
accessible and usable in the most democratic ways possible.

Therefore, for all our current and future software offerings, we'll perform a constant and well-educated evaluation whether
we can safely open source them in parts or even completely, as long as this appears to be non-harmful towards achieving
the project's main goal.

Also, we're constantly striving to keep our software offerings as accessible as possible when it comes to services which
cannot be run or managed by everyone - For example our Harmony Speech TTS Engine. As long as this project exists,
we'll be trying out utmost to provide free tiers for personal and public research use of our software and APIs.

However, at the same time we'll also ensure everyone who supports us or actively joins forces with us on our journey, gets
something proper back in turn. Therefore we're also maintaining a Patreon Page with different supporter tiers, as we are
open towards collaboration with other businesses.

### How to reach out to us

#### If you want to collaborate or support this Project financially:

Feel free to join our Discord Server and / or subscribe to our Patreon - Even $1 helps us drive this project forward.

![Harmony.AI Discord Server](docs/images/discord32.png) [Harmony.AI Discord Server](https://discord.gg/f6RQyhNPX8)

![Harmony.AI Discord Server](docs/images/patreon32.png) [Harmony.AI Patreon](https://patreon.com/harmony_ai)

#### If you want to use our software commercially or discuss a business or development partnership:

Contact us directly via: [contact@project-harmony.ai](mailto:contact@project-harmony.ai)

---
&copy; 2024 Harmony AI Solutions & Contributors

Licensed under the Apache 2.0 License
