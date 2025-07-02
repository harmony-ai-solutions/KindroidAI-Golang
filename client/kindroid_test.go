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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type KindroidAITestSuite struct {
	suite.Suite
	Server *httptest.Server
	Client *KindroidAI
}

func (suite *KindroidAITestSuite) SetupTest() {

	suite.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/send-message":
			// Validate request method
			suite.Equal(http.MethodPost, r.Method, "Expected method POST")

			// Validate headers
			authHeader := r.Header.Get("Authorization")
			suite.Equal("Bearer test_api_key", authHeader, "Invalid Authorization header")

			// Validate request body for SendMessage and SendMessageAdvanced
			bodyBytes, _ := io.ReadAll(r.Body)

			// For SendMessage (basic)
			if string(bodyBytes) == `{"ai_id":"test_ai_id","message":"Hello","stream":false}` {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"Hello, user!"`))
				return
			}

			// For SendMessageAdvanced (with multimedia options)
			expectedAdvancedBody := `{"ai_id":"test_ai_id","message":"Hello advanced","stream":true,"image_urls":["http://example.com/img.jpg"],"image_description":"a test image"}`
			if suite.JSONEq(expectedAdvancedBody, string(bodyBytes), "Invalid advanced request body") {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"Hello, advanced user!"`))
				return
			}

			w.WriteHeader(http.StatusBadRequest) // Fallback for unexpected body

		case "/chat-break":
			// Validate request method
			suite.Equal(http.MethodPost, r.Method, "Expected method POST")

			// Validate headers
			authHeader := r.Header.Get("Authorization")
			suite.Equal("Bearer test_api_key", authHeader, "Invalid Authorization header")

			// Validate request body
			bodyBytes, _ := io.ReadAll(r.Body)
			expectedBody := `{"ai_id":"test_ai_id","greeting":"Hello again"}`
			suite.JSONEq(expectedBody, string(bodyBytes), "Invalid request body")

			// Mock response
			w.WriteHeader(http.StatusOK)

		case "/check-user-subscription":
			suite.Equal(http.MethodPost, r.Method, "Expected method POST")
			authHeader := r.Header.Get("Authorization")
			suite.Equal("Bearer test_api_key", authHeader, "Invalid Authorization header")
			bodyBytes, _ := io.ReadAll(r.Body)
			suite.JSONEq("{}", string(bodyBytes), "Expected empty JSON body")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"uid":"test_uid","status":"OK","isSubscribedBase":true,"subscriptionPlatformBase":"web","gracePeriodBase":null,"isSubscribedAddon1":false,"subscriptionPlatformAddon1":null,"gracePeriodAddon1":null,"isSubscribedAddon2":false,"subscriptionPlatformAddon2":null,"gracePeriodAddon2":null}`))

		case "/audio-inference":
			suite.Equal(http.MethodPost, r.Method, "Expected method POST")
			authHeader := r.Header.Get("Authorization")
			suite.Equal("Bearer test_api_key", authHeader, "Invalid Authorization header")
			bodyBytes, _ := io.ReadAll(r.Body)
			expectedBody := `{"ai_id":"test_ai_id","messageID":"test_message_id"}`
			suite.JSONEq(expectedBody, string(bodyBytes), "Invalid request body")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`"OK"`)) // As per HAR, simple "OK" response

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	suite.Client = NewKindroidAI("test_api_key", "test_ai_id")
	suite.Client.BaseURL = suite.Server.URL
}

func (suite *KindroidAITestSuite) TearDownTest() {
	suite.Server.Close()
}

func (suite *KindroidAITestSuite) TestSendMessage() {
	response, err := suite.Client.SendMessage("Hello")
	suite.NoError(err, "SendMessage returned an error")
	suite.Equal(`"Hello, user!"`, response, "Unexpected response")
}

func (suite *KindroidAITestSuite) TestSendMessageAdvanced() {
	options := SendMessageOptions{
		AIID:             suite.Client.KindroidID,
		Message:          "Hello advanced",
		Stream:           true,
		ImageURLs:        []string{"http://example.com/img.jpg"},
		ImageDescription: func(s string) *string { return &s }("a test image"),
	}
	response, err := suite.Client.SendMessageAdvanced(options)
	suite.NoError(err, "SendMessageAdvanced returned an error")
	suite.Equal(`"Hello, advanced user!"`, response, "Unexpected advanced response")
}

func (suite *KindroidAITestSuite) TestChatBreak() {
	err := suite.Client.ChatBreak("Hello again")
	suite.NoError(err, "ChatBreak returned an error")
}

func (suite *KindroidAITestSuite) TestCheckUserSubscription() {
	subInfo, err := suite.Client.CheckUserSubscription()
	suite.NoError(err, "CheckUserSubscription returned an error")
	suite.NotNil(subInfo, "SubscriptionInfo should not be nil")
	suite.Equal("test_uid", subInfo.UID)
	suite.Equal("OK", subInfo.Status)
	suite.True(subInfo.IsSubscribedBase)
	suite.Equal("web", subInfo.SubscriptionPlatformBase)
}

func (suite *KindroidAITestSuite) TestAudioInference() {
	err := suite.Client.AudioInference("test_message_id")
	suite.NoError(err, "AudioInference returned an error")
}

func (suite *KindroidAITestSuite) TestExtractUserIDFromJWT() {
	// A dummy JWT with a "user_id" claim for testing purposes.
	// Header: {"alg":"HS256","typ":"JWT"}
	// Payload: {"user_id":"dummy_user_123","name":"Test User"}
	// Signature is not validated by ParseUnverified
	dummyJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZHVtbXlfdXNlcl8xMjMiLCJuYW1lIjoiVGVzdCBVc2VyIn0.signature_placeholder"

	// Create a temporary client with the dummy JWT
	tempClient := NewKindroidAI(dummyJWT, "any_ai_id")

	suite.Equal("dummy_user_123", tempClient.UserID, "UserID should be extracted correctly from JWT")

	// Test with a malformed JWT
	malformedJWT := "not.a.jwt"
	tempClientMalformed := NewKindroidAI(malformedJWT, "any_ai_id")
	suite.Empty(tempClientMalformed.UserID, "UserID should be empty for malformed JWT")

	// Test with JWT missing user_id claim
	jwtMissingClaim := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiVGVzdCBVc2VyIn0.signature_placeholder"
	tempClientMissingClaim := NewKindroidAI(jwtMissingClaim, "any_ai_id")
	suite.Empty(tempClientMissingClaim.UserID, "UserID should be empty for JWT missing user_id claim")
}

func TestKindroidAITestSuite(t *testing.T) {
	suite.Run(t, new(KindroidAITestSuite))
}
