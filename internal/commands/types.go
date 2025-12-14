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
	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/database"
)

// CommandContext holds all context for command execution
type CommandContext struct {
	Session     *discordgo.Session
	Message     *discordgo.MessageCreate
	Args        []string   // Arguments after command
	CleanArgs   []string   // Arguments with mentions resolved to usernames
	GuildConfig *database.GuildConfig
	Bot         interface{} // Will be *bot.Bot, using interface to avoid circular import
}

// CommandResponse represents the result of command execution
type CommandResponse struct {
	Content string                  // Plain text response
	Embed   *discordgo.MessageEmbed // Embed response
	File    *discordgo.File         // File attachment
	Files   []*discordgo.File       // Multiple file attachments
	Reply   bool                    // Whether to mention author in reply
}

// CommandProps defines command metadata
type CommandProps struct {
	Triggers        []string // Command triggers (first is primary)
	Description     string   // Command description
	Usage           string   // Usage pattern (e.g., "{command} <text>")
	Category        string   // Category name
	Cooldown        int64    // Cooldown in milliseconds (default: 3000)
	CooldownMessage string   // Custom cooldown message ({cooldown} is replaced with time)
	Permissions     []int64  // Required Discord permissions
	IsNSFW          bool     // NSFW flag
	OwnerOnly       bool     // Developer-only flag
	MissingArgs     string   // Message when args required but missing
}

// Command is the base interface all commands must implement
type Command interface {
	Props() CommandProps
	Run(ctx *CommandContext) (*CommandResponse, error)
}

// Helper functions for CommandContext
func (ctx *CommandContext) Reply(content string) *CommandResponse {
	return &CommandResponse{Content: content}
}

func (ctx *CommandContext) ReplyEmbed(embed *discordgo.MessageEmbed) *CommandResponse {
	return &CommandResponse{Embed: embed}
}

func (ctx *CommandContext) ReplyFile(name string, data []byte) *CommandResponse {
	return &CommandResponse{
		File: &discordgo.File{
			Name:        name,
			ContentType: "application/octet-stream",
			Reader:      bytesReader{data: data},
		},
	}
}

// Simple bytes reader for file responses
type bytesReader struct {
	data []byte
	pos  int
}

func (r bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, nil
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
