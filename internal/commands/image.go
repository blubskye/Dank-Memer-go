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
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ImageCommand handles image manipulation commands that use an external API
type ImageCommand struct {
	Properties   CommandProps
	RequestURL   string // URL with $URL placeholder for data
	TextOnly     bool   // If true, only accepts text input (no avatar)
	TextLimit    int    // Maximum text length
	RequiredArgs string // Error message when args are required but missing
	DoubleAvatar bool   // If true, uses two avatars (author and mentioned user)
	Format       string // Output format (png, gif, etc.)
}

func NewImageCommand(props CommandProps, requestURL string) *ImageCommand {
	return &ImageCommand{
		Properties: props,
		RequestURL: requestURL,
		Format:     "png",
	}
}

func (c *ImageCommand) Props() CommandProps {
	props := c.Properties

	// Apply defaults
	if props.Cooldown == 0 {
		props.Cooldown = 5000
	}
	if props.Category == "" {
		props.Category = "Image Manipulation"
	}

	// Add required permissions
	props.Permissions = append(props.Permissions,
		discordgo.PermissionEmbedLinks,
		discordgo.PermissionAttachFiles)

	return props
}

func (c *ImageCommand) Run(ctx *CommandContext) (*CommandResponse, error) {
	// Parse data source (avatar URL or text)
	dataSrc, err := c.parseDataSource(ctx)
	if err != nil {
		return &CommandResponse{Content: err.Error()}, nil
	}
	if dataSrc == "" {
		return nil, nil // Error already sent
	}

	// Get the bot interface
	bot, ok := ctx.Bot.(imageGenBot)
	if !ok {
		return nil, fmt.Errorf("bot does not support image generation")
	}

	// Make API request
	var imageData []byte
	if c.RequestURL != "" {
		url := strings.Replace(c.RequestURL, "$URL", dataSrc, 1)
		imageData, err = bot.GetImageGen().GenerateCustom("", map[string]string{
			"url": url,
		})
	} else {
		// Use default API endpoint
		imageData, err = bot.GetImageGen().Generate("/"+c.Properties.Triggers[0], dataSrc)
	}

	if err != nil {
		return &CommandResponse{Content: fmt.Sprintf("Error generating image: %s", err.Error())}, nil
	}

	format := c.Format
	if format == "" {
		format = "png"
	}

	return &CommandResponse{
		File: &discordgo.File{
			Name:   c.Properties.Triggers[0] + "." + format,
			Reader: bytes.NewReader(imageData),
		},
	}, nil
}

func (c *ImageCommand) parseDataSource(ctx *CommandContext) (string, error) {
	if c.TextOnly {
		// Text-only command
		if c.RequiredArgs != "" && len(ctx.Args) == 0 {
			return "", fmt.Errorf(c.RequiredArgs)
		}

		text := strings.Join(ctx.Args, " ")
		if c.TextLimit > 0 && len(text) > c.TextLimit {
			return "", fmt.Errorf("Too long. You're %d characters over the limit!", len(text)-c.TextLimit)
		}

		return text, nil
	}

	// Get avatar URL
	avatarURL := c.getAvatarURL(ctx)

	// Check for URL in args
	if len(ctx.Args) > 0 {
		argText := strings.Join(ctx.Args, " ")
		if isImageURL(argText) {
			avatarURL = strings.ReplaceAll(argText, "gif", "png")
			avatarURL = strings.ReplaceAll(avatarURL, "webp", "png")
		}
	}

	if c.RequiredArgs != "" && len(ctx.Args) == 0 {
		return "", fmt.Errorf(c.RequiredArgs)
	}

	if c.RequiredArgs != "" {
		// Command needs both avatar and text
		text := strings.Join(ctx.Args, " ")
		if c.TextLimit > 0 && len(text) > c.TextLimit {
			return "", fmt.Errorf("Too long. You're %d characters over the limit!", len(text)-c.TextLimit)
		}

		// Check for ASCII only
		if !isASCII(text) {
			return "", fmt.Errorf("Your argument contains invalid characters. Please try again.")
		}

		// Return JSON array of [avatar, text]
		data, _ := json.Marshal([]string{avatarURL, text})
		return string(data), nil
	}

	if c.DoubleAvatar {
		// Get second avatar (author or bot)
		var authorURL string
		if len(ctx.Message.Mentions) > 0 {
			authorURL = ctx.Message.Author.AvatarURL("1024")
		} else {
			authorURL = ctx.Session.State.User.AvatarURL("1024")
		}
		data, _ := json.Marshal([]string{avatarURL, authorURL})
		return string(data), nil
	}

	return avatarURL, nil
}

func (c *ImageCommand) getAvatarURL(ctx *CommandContext) string {
	// Check for mentions
	if len(ctx.Message.Mentions) > 0 {
		return ctx.Message.Mentions[0].AvatarURL("1024")
	}
	return ctx.Message.Author.AvatarURL("1024")
}

func isImageURL(url string) bool {
	lower := strings.ToLower(url)
	exts := []string{".jpg", ".jpeg", ".gif", ".png", ".webp"}
	for _, ext := range exts {
		if strings.Contains(lower, ext) {
			return true
		}
	}
	return false
}

func isASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

// Interface to access bot's image generator
type imageGenBot interface {
	GetImageGen() ImageGenerator
}

type ImageGenerator interface {
	Generate(endpoint, data string) ([]byte, error)
	GenerateCustom(endpoint string, params map[string]string) ([]byte, error)
}

// Mention parsing regex
var mentionRegex = regexp.MustCompile(`<@!?(\d+)>`)
