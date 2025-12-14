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
	"strings"

	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"google", "lmgtfy"},
			Description: "Sick of someone asking dumb questions? LMGTFY it for them!",
			Usage:       "{command} search terms",
			MissingArgs: "Hey, what do you want me to google?",
			Category:    "Fun Commands",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			if len(ctx.Args) == 0 {
				return &commands.CommandResponse{Content: "Hey, what do you want me to google?"}, nil
			}

			query := strings.Join(ctx.Args, "+")
			url := fmt.Sprintf("http://lmgtfy.com/?q=%s", query)

			return &commands.CommandResponse{Content: url}, nil
		},
	})
}
