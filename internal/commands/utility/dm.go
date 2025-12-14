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
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"dm", "slideintothedms"},
			Description: "melmsie stinks",
			Usage:       "{command} <id> <shit>",
			OwnerOnly:   true,
			Category:    "Utility Commands",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			if len(ctx.Args) < 2 {
				return &commands.CommandResponse{Content: "Usage: dm <user_id> <message>"}, nil
			}

			userID := ctx.Args[0]
			message := strings.Join(ctx.Args[1:], " ")

			// Create DM channel
			channel, err := ctx.Session.UserChannelCreate(userID)
			if err != nil {
				ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "❌")
				return &commands.CommandResponse{Content: "**Fuck!** *" + err.Error() + "*"}, nil
			}

			// Send message
			embed := &discordgo.MessageEmbed{
				Title:       "📫 You have received a message from the developers!",
				Description: message,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "To reply, please use pls vent.",
				},
				Color: utils.RandomColor(),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, embed)
			if err != nil {
				ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "❌")
				return &commands.CommandResponse{Content: "**Fuck!** *" + err.Error() + "*"}, nil
			}

			ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "📧")
			return nil, nil
		},
	})
}
