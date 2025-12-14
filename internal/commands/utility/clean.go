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
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"clean", "purge"},
			Description: "Will quickly clean the last 10 messages, or however many you specify.",
			Cooldown:    5000,
			Category:    "Utility Commands",
			Permissions: []int64{
				discordgo.PermissionManageMessages,
				discordgo.PermissionReadMessageHistory,
			},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Get the bot's user ID
			botUserID := ctx.Session.State.User.ID

			// Determine how many messages to clean
			count := 10
			if len(ctx.Args) > 0 {
				if num, err := strconv.Atoi(ctx.Args[0]); err == nil && num > 0 && num <= 100 {
					count = num
				}
			}

			// Get messages from channel
			messages, err := ctx.Session.ChannelMessages(ctx.Message.ChannelID, 100, "", "", "")
			if err != nil {
				return nil, err
			}

			// Filter messages by bot author and age (< 14 days)
			cutoff := time.Now().Add(-14 * 24 * time.Hour)
			var toDelete []string

			for _, msg := range messages {
				if msg.Author.ID == botUserID {
					msgTime := time.Time(msg.Timestamp)
					if msgTime.After(cutoff) {
						toDelete = append(toDelete, msg.ID)
						if len(toDelete) >= count {
							break
						}
					}
				}
			}

			if len(toDelete) == 0 {
				return nil, nil // No messages to delete
			}

			// Delete messages
			if len(toDelete) == 1 {
				ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, toDelete[0])
			} else {
				ctx.Session.ChannelMessagesBulkDelete(ctx.Message.ChannelID, toDelete)
			}

			return nil, nil
		},
	})
}
