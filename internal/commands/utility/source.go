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
			Triggers:    []string{"source", "sourcecode", "github", "code"},
			Description: "Get the source code for this bot (AGPL-3.0 licensed)",
			Category:    "Utility Commands",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			embed := &discordgo.MessageEmbed{
				Title:       "Dank Memer Source Code",
				Description: "This bot is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.\n\nYou have the right to view, modify, and redistribute this source code under the same license terms.",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Source Code",
						Value:  "[GitHub Repository](https://github.com/blubskye/Dank-Memer-go)",
						Inline: true,
					},
					{
						Name:   "License",
						Value:  "[AGPL-3.0](https://www.gnu.org/licenses/agpl-3.0.html)",
						Inline: true,
					},
					{
						Name:   "Original Project",
						Value:  "[Dank-Memer (Node.js)](https://github.com/melmsie/Dank-Memer)",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Free software means you have the freedom to run, study, share, and modify the software.",
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
