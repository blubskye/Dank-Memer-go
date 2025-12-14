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

package currency

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"coins", "balance", "bal"},
			Description: "u got dis many coins ok",
			Category:    "Currency",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Get database from bot
			b, ok := ctx.Bot.(dbBot)
			if !ok {
				return nil, fmt.Errorf("cannot access database")
			}

			// Get balance
			coins, err := b.GetDB().GetCoins(ctx.Message.Author.ID)
			if err != nil {
				return nil, err
			}

			embed := &discordgo.MessageEmbed{
				Title:       "how many coins you got fam?",
				Description: fmt.Sprintf("oh okay u got this many: %d", coins),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://dankmemer.lol/coin.png",
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "dont spend it all in one place ok",
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
