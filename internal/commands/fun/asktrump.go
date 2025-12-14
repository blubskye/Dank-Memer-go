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

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"asktrump", "donald"},
			Description: "Ask the president whatever you'd like!",
			Usage:       "{command} <question>",
			Category:    "Fun Commands",
			MissingArgs: "You gotta give me something to ask Trump :eyes:",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			question := strings.Join(ctx.Args, " ")

			// Count question marks for exclamation points
			qCount := strings.Count(question, "?")
			exclamations := strings.Repeat("!", qCount)

			response := strings.ToUpper(utils.GetTrumpResponse()) + exclamations
			photo := utils.GetTrumpPhoto()

			embed := &discordgo.MessageEmbed{
				Description: fmt.Sprintf("\n%s: %s\n\nTrump: %s",
					ctx.Message.Author.Username, question, response),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: photo,
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
