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
	"encoding/json"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/external"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"chucknorris", "chuck", "norris"},
			Description: "Let's learn about God",
			Category:    "Fun Commands",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			data, err := external.Get("http://api.icndb.com/jokes/random", nil)
			if err != nil {
				return &commands.CommandResponse{Content: "Couldn't fetch a Chuck Norris joke right now"}, nil
			}

			var result struct {
				Value struct {
					Joke string `json:"joke"`
				} `json:"value"`
			}
			if err := json.Unmarshal(data, &result); err != nil {
				return &commands.CommandResponse{Content: "Couldn't parse the joke"}, nil
			}

			joke := strings.ReplaceAll(result.Value.Joke, "&quot;", "\"")

			embed := &discordgo.MessageEmbed{
				Title:       "👊 Chuck Norris 👊",
				Description: joke,
				Color:       utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
