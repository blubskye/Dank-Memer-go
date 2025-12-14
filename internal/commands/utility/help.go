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
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"help", "cmds", "commands"},
			Description: "See a list of commands available.",
			Category:    "Utility",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			b, ok := ctx.Bot.(helpBot)
			if !ok {
				return nil, fmt.Errorf("cannot access command registry")
			}

			// If no args, show all commands by category
			if len(ctx.Args) == 0 {
				allCmds := b.GetCommands().GetAll()

				// Group by category
				categories := make(map[string][]string)
				for _, cmd := range allCmds {
					props := cmd.Props()
					if props.OwnerOnly {
						continue
					}

					category := props.Category
					if category == "" {
						category = "Other"
					}

					categories[category] = append(categories[category], props.Triggers[0])
				}

				// Sort categories
				var catNames []string
				for name := range categories {
					catNames = append(catNames, name)
				}
				sort.Strings(catNames)

				// Build fields
				var fields []*discordgo.MessageEmbedField
				for _, catName := range catNames {
					cmds := categories[catName]
					sort.Strings(cmds)
					fields = append(fields, &discordgo.MessageEmbedField{
						Name:  catName,
						Value: strings.Join(cmds, ", "),
					})
				}

				embed := &discordgo.MessageEmbed{
					Title:       "Available Commands",
					Description: "Auto posting memes, shorter cooldowns, custom commands and more coming on the premium bot later this week. Use pls patreon to see how to get it!",
					Fields:      fields,
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Hello darkness my old friend...",
					},
					Color: utils.RandomColor(),
				}

				return &commands.CommandResponse{Embed: embed}, nil
			}

			// Show specific command info
			cmdName := strings.ToLower(ctx.Args[0])
			cmd := b.GetCommands().Find(cmdName)
			if cmd == nil {
				return &commands.CommandResponse{Content: "Command not found."}, nil
			}

			props := cmd.Props()
			prefix := ctx.GuildConfig.Prefix

			usage := props.Usage
			if usage == "" {
				usage = "{command}"
			}
			usage = strings.ReplaceAll(usage, "{command}", prefix+" "+props.Triggers[0])

			embed := &discordgo.MessageEmbed{
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Description:", Value: props.Description},
					{Name: "Usage:", Value: "```" + usage + "```"},
					{Name: "Triggers:", Value: strings.Join(props.Triggers, ", ")},
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}

type helpBot interface {
	GetCommands() *commands.Registry
}
