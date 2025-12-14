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

package voice

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"stop"},
			Description: "STOP FARTING",
			Category:    "Voice Commands",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Check if user is in a voice channel
			voiceState := getUserVoiceState(ctx)
			if voiceState == nil || voiceState.ChannelID == "" {
				return &commands.CommandResponse{Content: "join a voice channel fam"}, nil
			}

			// Get the bot interface
			botInterface, ok := ctx.Bot.(interface {
				IsVoicePlaying(guildID string) bool
				StopAudio(guildID string) error
			})
			if !ok {
				return &commands.CommandResponse{Content: "Voice not supported"}, nil
			}

			// Check if playing
			if !botInterface.IsVoicePlaying(ctx.Message.GuildID) {
				return &commands.CommandResponse{Content: "I'm not playing anything right now!"}, nil
			}

			// Stop playback
			err := botInterface.StopAudio(ctx.Message.GuildID)
			if err != nil {
				return &commands.CommandResponse{Content: "Error stopping playback"}, nil
			}

			// Add reaction
			ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "❌")

			return nil, nil
		},
	})
}

func getUserVoiceState(ctx *commands.CommandContext) *discordgo.VoiceState {
	// Try to get from state cache
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err == nil {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == ctx.Message.Author.ID {
				return vs
			}
		}
	}

	// Fallback: fetch guild
	guild, err = ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			return vs
		}
	}

	return nil
}
