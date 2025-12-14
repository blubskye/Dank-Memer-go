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

package utility

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"invite", "support", "server"},
			Description: "Get an invite for the bot or to the support server.",
			Category:    "Utility Commands",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Try to get URLs from config
			inviteURL := "https://goo.gl/BPWvB9"
			supportURL := "https://discord.gg/ebUqc7F"

			botInterface, ok := ctx.Bot.(interface {
				GetConfig() *utils.Config
			})
			if ok {
				config := botInterface.GetConfig()
				if config.URLs.Invite != "" {
					inviteURL = config.URLs.Invite
				}
				if config.URLs.Support != "" {
					supportURL = config.URLs.Support
				}
			}

			embed := &discordgo.MessageEmbed{
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Add Dank Memer",
						Value:  "\n[Here](" + inviteURL + ")",
						Inline: true,
					},
					{
						Name:   "Join a Dank Server",
						Value:  "\n[Here](" + supportURL + ")",
						Inline: true,
					},
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
