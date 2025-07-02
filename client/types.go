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

// SubscriptionInfo represents the response structure for user subscription details.
type SubscriptionInfo struct {
	UID                        string  `json:"uid"`
	Status                     string  `json:"status"`
	IsSubscribedBase           bool    `json:"isSubscribedBase"`
	SubscriptionPlatformBase   string  `json:"subscriptionPlatformBase"`
	GracePeriodBase            *int    `json:"gracePeriodBase"` // Use pointer for nullable int
	IsSubscribedAddon1         bool    `json:"isSubscribedAddon1"`
	SubscriptionPlatformAddon1 *string `json:"subscriptionPlatformAddon1"` // Use pointer for nullable string
	GracePeriodAddon1          *int    `json:"gracePeriodAddon1"`
	IsSubscribedAddon2         bool    `json:"isSubscribedAddon2"`
	SubscriptionPlatformAddon2 *string `json:"subscriptionPlatformAddon2"`
	GracePeriodAddon2          *int    `json:"gracePeriodAddon2"`
}

// SendMessageOptions represents the request body for the enhanced SendMessage API.
type SendMessageOptions struct {
	AIID             string   `json:"ai_id"`
	Message          string   `json:"message"`
	Stream           bool     `json:"stream"`
	ImageURLs        []string `json:"image_urls,omitempty"`
	ImageDescription *string  `json:"image_description,omitempty"`
	VideoURL         *string  `json:"video_url,omitempty"`
	VideoDescription *string  `json:"video_description,omitempty"`
	InternetResponse *string  `json:"internet_response,omitempty"`
	LinkURL          *string  `json:"link_url,omitempty"`
	LinkDescription  *string  `json:"link_description,omitempty"`
}

// AudioInferenceRequest represents the request body for the Audio Inference API.
type AudioInferenceRequest struct {
	AIID      string `json:"ai_id"`
	MessageID string `json:"messageID"`
}
