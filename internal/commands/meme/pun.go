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

package meme

import (
	"encoding/json"

	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/external"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"pun", "dadjoke"},
			Description: "Are they dad jokes, or are they puns? Is there even a difference?",
			Category:    "Memey Commands",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			data, err := external.Get("https://icanhazdadjoke.com/", map[string]string{
				"Accept": "application/json",
			})
			if err != nil {
				return &commands.CommandResponse{Content: "Couldn't fetch a dad joke right now"}, nil
			}

			var result struct {
				Joke string `json:"joke"`
			}
			if err := json.Unmarshal(data, &result); err != nil {
				return &commands.CommandResponse{Content: "Couldn't parse the joke"}, nil
			}

			return &commands.CommandResponse{Content: result.Joke}, nil
		},
	})
}
