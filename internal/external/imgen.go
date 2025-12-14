// Dank Memer - A Discord bot
// Copyright (C) 2025 Dank Memer
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package external

import (
	"fmt"
	"net/url"
	"strings"
)

// ImageGenClient handles requests to the image generation API
type ImageGenClient struct {
	baseURL string
	apiKey  string
}

func NewImageGenClient(baseURL, apiKey string) *ImageGenClient {
	return &ImageGenClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
	}
}

// Generate makes a request to the image generation API
// endpoint should be like "/magik" or "/ban"
// data is the URL or text to process
func (c *ImageGenClient) Generate(endpoint, data string) ([]byte, error) {
	// Build URL with query parameter
	u := fmt.Sprintf("%s%s?avatar1=%s", c.baseURL, endpoint, url.QueryEscape(data))

	headers := map[string]string{}
	if c.apiKey != "" {
		headers["Authorization"] = c.apiKey
	}

	return Get(u, headers)
}

// GenerateWithText makes a request with text data
func (c *ImageGenClient) GenerateWithText(endpoint, text string) ([]byte, error) {
	u := fmt.Sprintf("%s%s?text=%s", c.baseURL, endpoint, url.QueryEscape(text))

	headers := map[string]string{}
	if c.apiKey != "" {
		headers["Authorization"] = c.apiKey
	}

	return Get(u, headers)
}

// GenerateDouble makes a request with two avatar URLs
func (c *ImageGenClient) GenerateDouble(endpoint, avatar1, avatar2 string) ([]byte, error) {
	u := fmt.Sprintf("%s%s?avatar1=%s&avatar2=%s",
		c.baseURL, endpoint,
		url.QueryEscape(avatar1),
		url.QueryEscape(avatar2))

	headers := map[string]string{}
	if c.apiKey != "" {
		headers["Authorization"] = c.apiKey
	}

	return Get(u, headers)
}

// GenerateCustom makes a request with custom query parameters
func (c *ImageGenClient) GenerateCustom(endpoint string, params map[string]string) ([]byte, error) {
	u := fmt.Sprintf("%s%s?", c.baseURL, endpoint)

	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	u += query.Encode()

	headers := map[string]string{}
	if c.apiKey != "" {
		headers["Authorization"] = c.apiKey
	}

	return Get(u, headers)
}
