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
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"mock"},
			Description: "Mock the stupid shit your friend says!",
			Usage:       "{command} <text to be mocked>",
			Category:    "Fun Commands",
			MissingArgs: "You gotta give me something to mock :eyes:",
			Permissions: []int64{discordgo.PermissionAttachFiles},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			text := strings.Join(ctx.Args, " ")

			// Replace c with k, v with c
			text = strings.ReplaceAll(text, "c", "k")
			text = strings.ReplaceAll(text, "C", "K")
			text = strings.ReplaceAll(text, "v", "c")
			text = strings.ReplaceAll(text, "V", "C")

			// Alternate case
			runes := []rune(text)
			for i, r := range runes {
				if i%2 == 1 {
					runes[i] = unicode.ToUpper(r)
				} else {
					runes[i] = unicode.ToLower(r)
				}
			}

			return &commands.CommandResponse{Content: string(runes)}, nil
		},
	})
}
