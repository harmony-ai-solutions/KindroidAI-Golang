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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// KindroidAI stores session parameters for interacting with the KindroidAI API.
type KindroidAI struct {
	APIKey     string
	KindroidID string
	BaseURL    string
	Client     *http.Client
}

// NewKindroidAI initializes a new KindroidAI client.
func NewKindroidAI(apiKey, kindroidID string) *KindroidAI {
	return &KindroidAI{
		APIKey:     apiKey,
		KindroidID: kindroidID,
		BaseURL:    "https://api.kindroid.ai/v1",
		Client:     &http.Client{},
	}
}

// SendMessage sends a message to the AI and returns the response.
func (k *KindroidAI) SendMessage(message string) (string, error) {
	url := fmt.Sprintf("%s/send-message", k.BaseURL)
	requestBody := map[string]string{
		"ai_id":   k.KindroidID,
		"message": message,
	}
	jsonData, err := json.Marshal(requestBody)
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
