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

var developers = []string{
	"Melmsie#0001",
	"Aetheryx#2222",
	"CyberRonin#5517",
}

var contributors = []string{
	"Kromatic#0420",
}

var staff = []string{
	"Kromatic",
	"Lizard",
	"Squidaddy",
	"xXBuilderBXx",
	"Akko",
}

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"credits", "helpers"},
			Description: "Thanks to all of you!",
			Category:    "Utility Commands",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			devsStr := ""
			for _, dev := range developers {
				devsStr += dev + "\n"
			}

			contribStr := ""
			for _, contrib := range contributors {
				contribStr += contrib + "\n"
			}

			staffStr := ""
			for _, s := range staff {
				staffStr += s + "\n"
			}

			embed := &discordgo.MessageEmbed{
				Title: "Dank Memer Credits",
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Developers", Value: devsStr, Inline: true},
					{Name: "Contributors", Value: contribStr, Inline: true},
					{Name: "Support Server Staff", Value: staffStr, Inline: true},
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
