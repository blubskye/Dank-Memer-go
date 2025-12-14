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

package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/external"
	"github.com/dankmemer/bot/internal/utils"
)

// MediaCommand handles commands that fetch media from external APIs
type MediaCommand struct {
	Properties CommandProps
	RequestURL string            // URL to fetch from
	JSONKey    string            // Key to extract from JSON response (if empty, uses full response)
	TokenKey   string            // Config key for API token (if empty, no auth)
	Title      string            // Embed title
	Message    string            // Footer message
	PrependURL string            // URL prefix for relative URLs
	Headers    map[string]string // Additional headers
}

func NewMediaCommand(props CommandProps, requestURL string) *MediaCommand {
	return &MediaCommand{
		Properties: props,
		RequestURL: requestURL,
	}
}

func (c *MediaCommand) Props() CommandProps {
	props := c.Properties

	// Apply defaults
	if props.Cooldown == 0 {
		props.Cooldown = 3000
	}

	// Add required permissions
	props.Permissions = append(props.Permissions, discordgo.PermissionEmbedLinks)

	return props
}

func (c *MediaCommand) Run(ctx *CommandContext) (*CommandResponse, error) {
	// Build headers
	headers := make(map[string]string)
	for k, v := range c.Headers {
		headers[k] = v
	}

	// Get token from config if specified
	if c.TokenKey != "" {
		bot, ok := ctx.Bot.(mediaConfigBot)
		if ok {
			token := bot.GetConfigValue(c.TokenKey)
			if token != "" {
				headers["Authorization"] = token
				headers["Key"] = token
			}
		}
	}

	// Make request
	var imageURL string
	var err error

	if c.JSONKey != "" {
		// Parse JSON response
		var result map[string]interface{}
		err = external.GetJSON(c.RequestURL, headers, &result)
		if err != nil {
			return &CommandResponse{Content: fmt.Sprintf("Error fetching media: %s", err.Error())}, nil
		}

		// Extract value from JSON
		if val, ok := result[c.JSONKey]; ok {
			imageURL = fmt.Sprintf("%v", val)
		} else {
			return &CommandResponse{Content: "Could not find media in response"}, nil
		}
	} else {
		// Use raw response as URL
		data, err := external.Get(c.RequestURL, headers)
		if err != nil {
			return &CommandResponse{Content: fmt.Sprintf("Error fetching media: %s", err.Error())}, nil
		}

		// Try to parse as JSON first
		var jsonResponse map[string]interface{}
		if json.Unmarshal(data, &jsonResponse) == nil {
			// If it's JSON, try common keys
			for _, key := range []string{"url", "link", "image", "file"} {
				if val, ok := jsonResponse[key]; ok {
					imageURL = fmt.Sprintf("%v", val)
					break
				}
			}
		}
		if imageURL == "" {
			imageURL = strings.TrimSpace(string(data))
		}
	}

	// Skip mp4 videos (retry)
	if strings.HasSuffix(strings.ToLower(imageURL), ".mp4") {
		return c.Run(ctx) // Retry
	}

	// Prepend URL if needed
	if c.PrependURL != "" && !strings.HasPrefix(imageURL, "http") {
		imageURL = c.PrependURL + imageURL
	}

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title: c.Title,
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
		Color: utils.RandomColor(),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", ctx.Message.Author.Username),
		},
	}

	if c.Message != "" {
		embed.Footer.Text += " | " + c.Message
	}

	return &CommandResponse{Embed: embed}, nil
}

// Interface to access bot's config
type mediaConfigBot interface {
	GetConfigValue(key string) string
}
