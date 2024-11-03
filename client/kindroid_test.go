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

			// Validate request body
			bodyBytes, _ := io.ReadAll(r.Body)
			expectedBody := `{"ai_id":"test_ai_id","message":"Hello"}`
			suite.JSONEq(expectedBody, string(bodyBytes), "Invalid request body")

			// Mock response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`"Hello, user!"`))

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

func (suite *KindroidAITestSuite) TestChatBreak() {
	err := suite.Client.ChatBreak("Hello again")
	suite.NoError(err, "ChatBreak returned an error")
}

func TestKindroidAITestSuite(t *testing.T) {
	suite.Run(t, new(KindroidAITestSuite))
}
