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

package bot

import (
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/database"
	"github.com/dankmemer/bot/internal/external"
	"github.com/dankmemer/bot/internal/utils"
)

// RegisterCommands is called from main to register all commands
// This function is populated by command package init() functions
func RegisterCommands(b *Bot) {
	// Commands register themselves via the global registry
	// We just need to copy them to the bot's registry
	for _, cmd := range globalRegistry.GetAll() {
		b.Commands.Register(cmd)
	}
}

// Global registry that command packages use to register
var globalRegistry = commands.NewRegistry()

// Register is called by command packages in their init() functions
func Register(cmd commands.Command) {
	globalRegistry.Register(cmd)
}

// Implement interfaces required by command types

// GetDB returns the database instance
func (b *Bot) GetDB() *database.Database {
	return b.DB
}

// GetConfig returns the bot config
func (b *Bot) GetConfig() *utils.Config {
	return b.Config
}

// GetCommands returns the command registry
func (b *Bot) GetCommands() *commands.Registry {
	return b.Commands
}

// GetImageGen returns the image generation client
func (b *Bot) GetImageGen() *external.ImageGenClient {
	return b.ImageGen
}

// GetRedditClient returns the Reddit client
func (b *Bot) GetRedditClient() RedditClientAdapter {
	return &redditClientAdapter{b.RedditClient}
}

// GetConfigValue returns a config value by key
func (b *Bot) GetConfigValue(key string) string {
	switch key {
	case "mashape":
		return "" // Add to config if needed
	case "novo":
		return ""
	default:
		return ""
	}
}

// IsVoicePlaying checks if the bot is playing audio in a guild
func (b *Bot) IsVoicePlaying(guildID string) bool {
	// Check voice manager
	if b.VoiceManager != nil {
		return b.VoiceManager.IsPlaying(guildID)
	}
	return false
}

// PlayAudio plays an audio file in a voice channel
func (b *Bot) PlayAudio(guildID, channelID, audioPath string) error {
	if b.VoiceManager == nil {
		return nil
	}
	return b.VoiceManager.PlayAudio(guildID, channelID, audioPath)
}

// StopAudio stops audio playback in a guild
func (b *Bot) StopAudio(guildID string) error {
	if b.VoiceManager == nil {
		return nil
	}
	b.VoiceManager.Stop(guildID)
	return nil
}

// Reddit client adapter to match the interface
type RedditClientAdapter interface {
	FetchPosts(endpoint string) ([]redditPost, error)
}

type redditClientAdapter struct {
	client *external.RedditClient
}

type redditPost struct {
	Title     string
	URL       string
	Permalink string
	Author    string
	Selftext  string
	PostHint  string
}

func (a *redditClientAdapter) FetchPosts(endpoint string) ([]redditPost, error) {
	posts, err := a.client.FetchPosts(endpoint)
	if err != nil {
		return nil, err
	}

	result := make([]redditPost, len(posts))
	for i, p := range posts {
		result[i] = redditPost{
			Title:     p.Title,
			URL:       p.URL,
			Permalink: p.Permalink,
			Author:    p.Author,
			Selftext:  p.Selftext,
			PostHint:  p.PostHint,
		}
	}
	return result, nil
}
