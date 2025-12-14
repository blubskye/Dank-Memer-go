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

package animal

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

var redpandaImages = []string{
	"https://i.imgur.com/0K3F3R9.jpg",
	"https://i.imgur.com/wS7aJW2.jpg",
	"https://i.imgur.com/F6R5RlG.jpg",
	"https://i.imgur.com/8Jl0Y9O.jpg",
	"https://i.imgur.com/4wh25QM.jpg",
}

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"redpanda", "redboi"},
			Description: "See some cute red pandas!",
			Category:    "Animal Commands",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			embed := &discordgo.MessageEmbed{
				Title: "dawwwwwwww 🐼",
				Image: &discordgo.MessageEmbedImage{
					URL: utils.RandomInArray(redpandaImages),
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Requested by %s", ctx.Message.Author.Username),
				},
				Color: utils.RandomColor(),
			}
			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
