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

package fun

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"vent", "bother", "message"},
			Description: "I know I am your only friend, this should give you someone to vent to",
			Cooldown:    15000,
			MissingArgs: "What do you want to vent to me about?",
			Category:    "Fun Commands",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			if len(ctx.Args) == 0 {
				return &commands.CommandResponse{Content: "What do you want to vent to me about?"}, nil
			}

			// Get config for vent channel
			botInterface, ok := ctx.Bot.(interface {
				GetConfig() *utils.Config
			})
			if !ok {
				return &commands.CommandResponse{Content: "My owner is listening to your venting now, what a great listener. amirite ladies???? 😍"}, nil
			}

			config := botInterface.GetConfig()
			ventChannelID := config.VentChannel
			if ventChannelID == "" {
				return &commands.CommandResponse{Content: "My owner is listening to your venting now, what a great listener. amirite ladies???? 😍"}, nil
			}

			// Get channel and guild info
			channel, err := ctx.Session.State.Channel(ctx.Message.ChannelID)
			if err != nil {
				channel, _ = ctx.Session.Channel(ctx.Message.ChannelID)
			}

			var channelName, guildName, guildID string
			if channel != nil {
				channelName = channel.Name
				guild, err := ctx.Session.State.Guild(channel.GuildID)
				if err != nil {
					guild, _ = ctx.Session.Guild(channel.GuildID)
				}
				if guild != nil {
					guildName = guild.Name
					guildID = guild.ID
				}
			}

			// Send vent to channel
			embed := &discordgo.MessageEmbed{
				Title: "New Vent:",
				Author: &discordgo.MessageEmbedAuthor{
					Name: fmt.Sprintf("%s#%s | %s",
						ctx.Message.Author.Username,
						ctx.Message.Author.Discriminator,
						ctx.Message.Author.ID),
				},
				Description: strings.Join(ctx.Args, " "),
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Sent from:",
						Value: fmt.Sprintf("#%s in %s", channelName, guildName),
					},
				},
				Color:     utils.RandomColor(),
				Timestamp: time.Now().Format(time.RFC3339),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Guild ID: %s", guildID),
				},
			}

			ctx.Session.ChannelMessageSendEmbed(ventChannelID, embed)

			return &commands.CommandResponse{
				Content: "My owner is listening to your venting now, what a great listener. amirite ladies???? 😍",
			}, nil
		},
	})
}
