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
			Triggers:        []string{"daily"},
			Description:     "u got dis many coins ok",
			Category:        "Currency",
			Cooldown:        86400000, // 24 hours
			CooldownMessage: "I'm not made of money dude, wait {cooldown}",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Get database from bot
			b, ok := ctx.Bot.(dbBot)
			if !ok {
				return nil, fmt.Errorf("cannot access database")
			}

			// Add coins
			if err := b.GetDB().AddCoins(ctx.Message.Author.ID, 100); err != nil {
				return nil, err
			}

			// Get new balance
			coins, err := b.GetDB().GetCoins(ctx.Message.Author.ID)
			if err != nil {
				return nil, err
			}

			embed := &discordgo.MessageEmbed{
				Title:       "here are ur daily coins ok",
				Description: fmt.Sprintf("u got 100, now u have %d", coins),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://dankmemer.lol/coin.png",
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "spend it all in one place ok",
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}

// Interface to access bot's database
type dbBot interface {
	GetDB() database
}

type database interface {
	AddCoins(userID string, amount int64) error
	GetCoins(userID string) (int64, error)
}
