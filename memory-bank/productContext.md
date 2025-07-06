# Product Context: KindroidAI Golang Client

## Why this project exists
This project exists to provide a native Golang client for the KindroidAI API. It aims to simplify the integration of KindroidAI functionalities into Golang applications, making it easier for developers to build AI-driven features.

## Problems it solves
- **Lack of Official Golang Client**: Fills the gap for Golang developers who want to interact with the KindroidAI API without having to implement the HTTP requests and response parsing from scratch.
- **Simplified API Interaction**: Abstracts away the complexities of direct HTTP calls, authentication, and JSON serialization/deserialization, providing a clean and idiomatic Golang interface.
- **Accelerated Development**: Enables faster development of applications that leverage KindroidAI by providing a ready-to-use client.

## How it should work
The client should provide a straightforward API that mirrors the KindroidAI public API documentation. Developers should be able to:
- Initialize the client with an API key and Kindroid ID.
- Send messages to a Kindroid AI and receive responses.
- Reset chat sessions with a custom greeting.
- Handle API responses and errors gracefully.

## User experience goals (for developers using this client)
- **Ease of Use**: The client should be intuitive and easy to integrate into existing Golang projects.
- **Reliability**: API calls should be robust and handle network issues or API errors gracefully.
- **Clarity**: The client's methods and data structures should be clear and well-documented, reflecting the underlying KindroidAI API.
- **Efficiency**: Provide a performant way to interact with the KindroidAI API.
