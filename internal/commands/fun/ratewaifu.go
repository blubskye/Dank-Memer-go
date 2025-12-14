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
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

var waifuReactions = []string{
	"🤢", "😰", "😵", "😥", "😕", "😗", "😏", "😌", "☺", "💁",
}

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"waifu", "ratewaifu"},
			Description: "See how good of a waifu you are",
			Category:    "Fun Commands",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			var target string
			if len(ctx.Args) == 0 || strings.ToLower(ctx.Args[0]) == "me" {
				target = "You are a"
			} else if len(ctx.Message.Mentions) > 0 {
				target = ctx.Message.Mentions[0].Username + " is a"
			} else {
				target = strings.Join(ctx.Args, " ") + " is a"
			}

			rating := rand.Intn(100) + 1
			reactionIdx := rating / 10
			if reactionIdx > 9 {
				reactionIdx = 9
			}
			if reactionIdx < 0 {
				reactionIdx = 0
			}

			embed := &discordgo.MessageEmbed{
				Title:       "waifu r8 machine",
				Description: fmt.Sprintf("%s %d/100 waifu %s", target, rating, waifuReactions[reactionIdx]),
				Color:       utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
