# Active Context: KindroidAI Golang Client

## Current Work Focus
The primary focus has been on implementing a new, undocumented feature: fetching and decrypting chat history directly from Google's Firestore, mirroring the behavior of the Kindroid web application.

## Recent Changes
- **Added `GetChatHistory` Method**: Implemented a new method in `client/kindroid.go` to retrieve chat messages.
- **Firestore Integration**: Added `cloud.google.com/go/firestore` as a dependency to interact with the Firestore API.
- **JWT Parsing**: The client now automatically parses the provided `APIKey` as a JWT to extract the `UserID`.
- **Fallback Authentication**: If the `APIKey` is not a valid JWT, the client now falls back to using the `KINDROID_USER_ID` environment variable.
- **Message Decryption**: Implemented message decryption using the `UserID` as the key. This handles the `!enc:` format found in Firestore messages.
- **Updated Dependencies**: Added `github.com/Luzifer/go-openssl/v4` for decryption and `github.com/golang-jwt/jwt/v5` for token parsing.
- **New Example**: Created `examples/example_chat_history.go` to demonstrate the new functionality.
- **Documentation**: Updated `README.md` to reflect the new feature and authentication logic.

## Next Steps
- Final review of all changes.
- Update the project's `progress.md` to log the completion of this feature.

## Active Decisions and Considerations
- **Authentication Flexibility**: The decision to support both JWT and static API key + UserID provides greater flexibility for different use cases.
- **Decryption Key**: The successful decryption confirms that the `UserID` is the correct key for message content, which is a critical piece of information for the client's functionality.
- **Error Handling**: Implemented graceful fallbacks and warnings for JWT parsing and decryption failures to improve user experience.

## Important Patterns and Preferences
- **Go Modules**: Project uses Go Modules for dependency management.
- **Environment Variables**: API keys and IDs are expected via environment variables.
- **Testify for Testing**: `github.com/stretchr/testify` is the preferred testing framework.
- **Standard HTTP Client**: Direct use of `net/http` for API interactions.
- **Firestore for History**: Chat history is not retrieved via a simple REST API call but through a gRPC connection to Firestore, requiring the official Google Cloud SDK.

## Learnings and Project Insights
- The project now has a significant feature that goes beyond the documented public API, making it much more powerful.
- The process revealed a complex authentication and data storage mechanism involving JWTs, Firestore, and custom encryption.
- The `UserID` is the central piece of information for both authentication and decryption of user-specific data.
- The client is now capable of providing a more complete interaction with the Kindroid platform.
