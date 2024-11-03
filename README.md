# 💬 KindroidAI - Golang Client by Project Harmony.AI
![Tag](https://img.shields.io/github/license/harmony-ai-solutions/CharacterAI-Golang)

An unofficial API Client for CharacterAI, written in Golang.

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
Detailed documentation and Apidocs coming soon

## 📙 Example
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
&copy; 2023 Harmony AI Solutions & Contributors

Licensed under the Apache 2.0 License