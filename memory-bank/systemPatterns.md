# System Patterns: KindroidAI Golang Client

## System Architecture
The KindroidAI Golang Client is designed as a lightweight wrapper around the KindroidAI REST API. It follows a client-server interaction model where the Golang application acts as a client making HTTP requests to the KindroidAI API endpoints.

## Key Technical Decisions
- **Direct HTTP Client Usage**: Utilizes Go's standard `net/http` package for making API requests, providing direct control over HTTP interactions.
- **JSON for Request/Response**: All communication with the KindroidAI API is handled using JSON payloads for both requests and responses.
- **Bearer Token Authentication**: API authentication is managed via a `Bearer` token in the `Authorization` header for each request.
- **Environment Variable Configuration**: API key and Kindroid ID are expected to be provided via environment variables, promoting secure and flexible configuration.
- **Test-Driven Development (TDD) / Comprehensive Testing**: The project includes a robust test suite using `github.com/stretchr/testify` and `httptest` to ensure the client's functionality and API interactions are correct.

## Design Patterns in Use
- **Client Pattern**: The `KindroidAI` struct acts as a client, encapsulating the API key, Kindroid ID, base URL, and an `http.Client` instance. This centralizes API interaction logic.
- **Dependency Injection (Implicit)**: The `http.Client` is part of the `KindroidAI` struct, allowing for easy mocking in tests by replacing the default client with a test server.
- **Error Handling**: Standard Golang error return patterns are used, where functions return a value and an `error` type, allowing callers to handle errors explicitly.

## Component Relationships
- **`KindroidAI` struct**: The core component, responsible for holding configuration and orchestrating API calls.
- **`http.Client`**: Used by the `KindroidAI` struct to perform actual HTTP requests.
- **`SendMessage` method**: Interacts with the `/send-message` endpoint.
- **`ChatBreak` method**: Interacts with the `/chat-break` endpoint.
- **`main` package (example.go)**: Demonstrates how to use the `client` package, acting as a consumer of the `KindroidAI` client.
- **`kindroid_test.go`**: Contains unit tests that mock the KindroidAI API server to verify client behavior.

## Critical Implementation Paths
- **API Request Construction**: Ensuring correct URL formatting, request body marshaling to JSON, and setting appropriate headers (Content-Type, Authorization).
- **API Response Handling**: Reading and parsing HTTP responses, checking status codes, and handling potential errors from the API or network.
- **Authentication Flow**: Correctly attaching the Bearer token to every outgoing request.
