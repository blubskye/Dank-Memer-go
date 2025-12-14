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
			Triggers:    []string{"patreon", "donate", "gibmonies", "pay", "donut", "plsdonut"},
			Description: "See how you can donate to the bot and gain access to donor features!",
			Category:    "Utility Commands",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			patreonURL := "https://www.patreon.com/dankmemerbot"

			botInterface, ok := ctx.Bot.(interface {
				GetConfig() *utils.Config
			})
			if ok {
				config := botInterface.GetConfig()
				if config.URLs.Patreon != "" {
					patreonURL = config.URLs.Patreon
				}
			}

			embed := &discordgo.MessageEmbed{
				Title:       "Donate to Dank Memer on Patreon!",
				Description: "Help Melmsie keep the bot alive by donating to help pay server costs!",
				URL:         patreonURL,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "pls donut for my bot, im running low on money",
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
