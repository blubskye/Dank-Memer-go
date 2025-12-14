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

package text

import (
	"strings"
	"unicode"

	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
)

var specialCodes = map[rune]string{
	'0': ":zero:",
	'1': ":one:",
	'2': ":two:",
	'3': ":three:",
	'4': ":four:",
	'5': ":five:",
	'6': ":six:",
	'7': ":seven:",
	'8': ":eight:",
	'9': ":nine:",
	'#': ":hash:",
	'*': ":asterisk:",
	'?': ":grey_question:",
	'!': ":grey_exclamation:",
	' ': "   ",
}

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"emojify"},
			Description: "Make the bot say whatever you want with emojis!",
			Usage:       "{command} <what you want the bot to say>",
			Category:    "Text Commands",
			MissingArgs: "What do you want me to put into emojis?",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			text := strings.ToLower(strings.Join(ctx.Args, " "))

			var result strings.Builder
			for _, letter := range text {
				if letter >= 'a' && letter <= 'z' {
					result.WriteString(":regional_indicator_")
					result.WriteRune(letter)
					result.WriteString(": ")
				} else if code, ok := specialCodes[letter]; ok {
					result.WriteString(code)
					result.WriteString(" ")
				} else if unicode.IsPrint(letter) {
					result.WriteRune(letter)
				}
			}

			return &commands.CommandResponse{Content: result.String()}, nil
		},
	})
}
