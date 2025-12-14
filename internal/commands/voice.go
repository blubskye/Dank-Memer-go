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
	"fmt"
	"math/rand"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

// VoiceCommand handles commands that play audio in voice channels
type VoiceCommand struct {
	Properties   CommandProps
	Directory    string // Audio directory name (under assets/audio/)
	FileCount    int    // Number of audio files (for random selection)
	SingleFile   string // Single file name (if not random)
	Extension    string // Audio extension (opus, dca, mp3)
	Reaction     string // Emoji reaction to add
	ExistingConn string // Message for existing connection
}

func NewVoiceCommand(props CommandProps, directory string, fileCount int) *VoiceCommand {
	return &VoiceCommand{
		Properties: props,
		Directory:  directory,
		FileCount:  fileCount,
		Extension:  "opus",
	}
}

func (c *VoiceCommand) Props() CommandProps {
	props := c.Properties

	// Apply defaults
	if props.Cooldown == 0 {
		props.Cooldown = 10000
	}
	if props.Category == "" {
		props.Category = "Voice Commands"
	}

	// Add required permissions
	props.Permissions = append(props.Permissions, discordgo.PermissionAddReactions)

	return props
}

func (c *VoiceCommand) Run(ctx *CommandContext) (*CommandResponse, error) {
	// Check if user is in a voice channel
	voiceState, err := c.getUserVoiceState(ctx)
	if err != nil || voiceState == nil || voiceState.ChannelID == "" {
		return &CommandResponse{Content: "join a voice channel fam"}, nil
	}

	// Get the bot interface
	bot, ok := ctx.Bot.(voiceBot)
	if !ok {
		return nil, fmt.Errorf("bot does not support voice")
	}

	// Check if already playing in this guild
	if bot.IsVoicePlaying(ctx.Message.GuildID) {
		msg := c.ExistingConn
		if msg == "" {
			msg = "I'm already playing something. Please wait until the current sound is done."
		}
		return &CommandResponse{Content: msg}, nil
	}

	// Check voice channel permissions
	perms, err := ctx.Session.State.UserChannelPermissions(ctx.Session.State.User.ID, voiceState.ChannelID)
	if err == nil {
		requiredPerms := int64(discordgo.PermissionVoiceConnect | discordgo.PermissionVoiceSpeak)
		if perms&requiredPerms != requiredPerms {
			return &CommandResponse{
				Content: "Make sure I have `connect` and `speak` permissions in the voice channel!\n\nHow to do that: https://i.imgur.com/ugplJJO.gif",
			}, nil
		}
	}

	// Determine file to play
	var filename string
	if c.SingleFile != "" {
		filename = c.SingleFile
	} else {
		fileNum := rand.Intn(c.FileCount) + 1
		filename = fmt.Sprintf("%d", fileNum)
	}

	ext := c.Extension
	if ext == "" {
		ext = "opus"
	}

	audioPath := filepath.Join("assets", "audio", c.Directory, filename+"."+ext)

	// Add reaction
	if c.Reaction != "" {
		ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, c.Reaction)
	}

	// Play audio
	err = bot.PlayAudio(ctx.Message.GuildID, voiceState.ChannelID, audioPath)
	if err != nil {
		return &CommandResponse{Content: fmt.Sprintf("Error playing audio: %s", err.Error())}, nil
	}

	return nil, nil
}

func (c *VoiceCommand) getUserVoiceState(ctx *CommandContext) (*discordgo.VoiceState, error) {
	// Try to get from state cache
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err == nil {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == ctx.Message.Author.ID {
				return vs, nil
			}
		}
	}

	// Fallback: fetch guild
	guild, err = ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil, err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			return vs, nil
		}
	}

	return nil, nil
}

// Interface to access bot's voice functionality
type voiceBot interface {
	IsVoicePlaying(guildID string) bool
	PlayAudio(guildID, channelID, audioPath string) error
	StopAudio(guildID string) error
}
