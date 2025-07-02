// Package client
/*
Copyright © 2024 Harmony AI Solutions & Contributors

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
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/Luzifer/go-openssl/v4"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

// KindroidAI stores session parameters for interacting with the KindroidAI API.
type KindroidAI struct {
	APIKey     string
	KindroidID string
	BaseURL    string
	Client     *http.Client
	UserID     string // Added to store the UserID extracted from the JWT
}

// NewKindroidAI initializes a new KindroidAI client.
// It will attempt to extract the UserID from the apiKey if it's a JWT.
// If not, it will fall back to the KINDROID_USER_ID environment variable.
func NewKindroidAI(apiKey, kindroidID string) *KindroidAI {
	k := &KindroidAI{
		APIKey:     apiKey,
		KindroidID: kindroidID,
		BaseURL:    "https://api.kindroid.ai/v1",
		Client:     &http.Client{},
	}

	// Attempt to extract UserID from the APIKey (JWT)
	userID, err := k.extractUserIDFromJWT()
	if err != nil {
		// If JWT parsing fails, assume it's a regular API key and look for env var
		fmt.Printf("Info: Could not extract UserID from APIKey (not a JWT?). Falling back to KINDROID_USER_ID env var.\n")
		userID = os.Getenv("KINDROID_USER_ID")
	}

	if userID == "" {
		fmt.Printf("Warning: UserID is not set. Chat history features will not work.\n")
	}

	k.UserID = userID
	return k
}

// extractUserIDFromJWT parses the JWT (APIKey) to get the user ID.
func (k *KindroidAI) extractUserIDFromJWT() (string, error) {
	tokenString := k.APIKey
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid JWT claims format")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found or not a string in JWT claims")
	}

	return userID, nil
}

// SendMessage sends a message to the AI and returns the response.
// This is the basic version for backwards compatibility.
func (k *KindroidAI) SendMessage(message string) (string, error) {
	options := SendMessageOptions{
		AIID:    k.KindroidID,
		Message: message,
		Stream:  false, // Default to false for basic SendMessage
	}
	return k.SendMessageAdvanced(options)
}

// SendMessageAdvanced sends a message to the AI with advanced options and returns the response.
// This method supports multimedia, streaming, and other advanced features.
func (k *KindroidAI) SendMessageAdvanced(options SendMessageOptions) (string, error) {
	url := fmt.Sprintf("%s/send-message", k.BaseURL)
	jsonData, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.APIKey)

	resp, err := k.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// ChatBreak ends the current chat session and starts a new one with a customizable greeting sent by the AI.
func (k *KindroidAI) ChatBreak(greeting string) error {
	url := fmt.Sprintf("%s/chat-break", k.BaseURL)
	requestBody := map[string]string{
		"ai_id":    k.KindroidID,
		"greeting": greeting,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.APIKey)

	resp, err := k.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return nil
}

// CheckUserSubscription retrieves user subscription information.
//
// WARNING: This method uses an undocumented API endpoint discovered through
// network analysis. It may change or be removed without notice.
// Use at your own risk in production environments.
func (k *KindroidAI) CheckUserSubscription() (*SubscriptionInfo, error) {
	url := fmt.Sprintf("%s/check-user-subscription", k.BaseURL)
	// The HAR file shows an empty JSON object as the request body.
	jsonData := []byte("{}")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.APIKey)

	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var subInfo SubscriptionInfo
	err = json.Unmarshal(bodyBytes, &subInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription info: %w", err)
	}

	return &subInfo, nil
}

// AudioInference sends an audio inference request to the AI.
//
// WARNING: This method uses an undocumented API endpoint discovered through
// network analysis. It may change or be removed without notice.
// Use at your own risk in production environments.
func (k *KindroidAI) AudioInference(messageID string) error {
	url := fmt.Sprintf("%s/audio-inference", k.BaseURL)
	requestBody := AudioInferenceRequest{
		AIID:      k.KindroidID,
		MessageID: messageID,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.APIKey)

	resp, err := k.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// The HAR shows a simple "OK" response, no complex parsing needed.
	return nil
}

// decryptMessage decrypts a message string that is prefixed with "!enc:".
func (k *KindroidAI) decryptMessage(encryptedMsg string) (string, error) {
	if !strings.HasPrefix(encryptedMsg, "!enc:") {
		// Not encrypted, return as is
		return encryptedMsg, nil
	}

	trimmedMsg := strings.TrimPrefix(encryptedMsg, "!enc:")

	o := openssl.New()
	decrypted, err := o.DecryptBytes(k.UserID, []byte(trimmedMsg), openssl.BytesToKeyMD5)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %w", err)
	}

	return string(decrypted), nil
}

// GetChatHistory retrieves the most recent chat messages for a given AI from Firestore.
func (k *KindroidAI) GetChatHistory(ctx context.Context, aiID string, limit int) ([]ChatMessage, error) {
	if k.UserID == "" {
		return nil, fmt.Errorf("user ID not available; ensure APIKey is a valid JWT")
	}

	// Use the APIKey as the bearer token for Firestore authentication.
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: k.APIKey})
	opts := []option.ClientOption{option.WithTokenSource(ts)}

	// Initialize the Firestore client.
	// Note: The project ID is hardcoded based on HAR file analysis.
	client, err := firestore.NewClientWithDatabase(ctx, "kindroid-ai", "(default)", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}
	defer client.Close()

	// Construct the path to the "ChatMessages" collection.
	parentPath := fmt.Sprintf("Users/%s/AIs/%s", k.UserID, aiID)

	// Build the query.
	query := client.Collection(parentPath+"/ChatMessages").
		OrderBy("timestamp", firestore.Desc).
		Limit(limit)

	// Execute the query.
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %w", err)
	}

	// Parse and decrypt the documents.
	var messages []ChatMessage
	for _, doc := range docs {
		var msg ChatMessage
		if err := doc.DataTo(&msg); err != nil {
			fmt.Printf("Warning: Failed to parse chat message document %s: %v\n", doc.Ref.ID, err)
			continue
		}
		msg.ID = doc.Ref.ID

		// Decrypt the message content if it's encrypted.
		decryptedText, err := k.decryptMessage(msg.Message)
		if err != nil {
			fmt.Printf("Warning: Failed to decrypt message for doc %s: %v\n", doc.Ref.ID, err)
			msg.Message = "[DECRYPTION FAILED]"
		} else {
			msg.Message = decryptedText
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
