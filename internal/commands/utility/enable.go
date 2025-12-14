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
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"enable"},
			Description: "Use this command to enable disabled commands.",
			Category:    "Utility Commands",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Check permissions
			botInterface, ok := ctx.Bot.(interface {
				GetConfig() *utils.Config
			})

			hasPermission := false

			// Check if user has manage guild permission
			perms, err := ctx.Session.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err == nil && perms&discordgo.PermissionManageServer != 0 {
				hasPermission = true
			}

			// Check if user is a dev
			if ok {
				config := botInterface.GetConfig()
				for _, dev := range config.Devs {
					if dev == ctx.Message.Author.ID {
						hasPermission = true
						break
					}
				}
			}

			if !hasPermission {
				return &commands.CommandResponse{
					Content: "You are not authorized to use this command. You must have `Manage Server` to enable commands.",
				}, nil
			}

			if len(ctx.Args) == 0 {
				return &commands.CommandResponse{
					Content: fmt.Sprintf("Specify a command to enable, or multiple.\n\nExample: `%s enable meme trigger shitsound` or `%s enable meme`",
						ctx.GuildConfig.Prefix, ctx.GuildConfig.Prefix),
				}, nil
			}

			// Get the command registry
			cmdBot, ok := ctx.Bot.(interface {
				GetCommands() *commands.Registry
			})
			if !ok {
				return &commands.CommandResponse{Content: "Error accessing commands"}, nil
			}

			registry := cmdBot.GetCommands()

			// Remove duplicates and normalize command names
			seen := make(map[string]bool)
			var normalizedArgs []string

			for _, arg := range ctx.Args {
				arg = strings.ToLower(arg)

				// Find the canonical command name
				cmd := registry.Find(arg)
				var cmdName string
				if cmd != nil {
					cmdName = cmd.Props().Triggers[0]
				} else if arg == "nsfw" {
					cmdName = "nsfw"
				} else {
					continue
				}

				if !seen[cmdName] {
					seen[cmdName] = true
					normalizedArgs = append(normalizedArgs, cmdName)
				}
			}

			if len(normalizedArgs) == 0 {
				return &commands.CommandResponse{
					Content: "No valid commands specified.",
				}, nil
			}

			// Check which are not disabled
			var notDisabled []string
			for _, cmdName := range normalizedArgs {
				isDisabled := false
				for _, disabled := range ctx.GuildConfig.DisabledCommands {
					if disabled == cmdName {
						isDisabled = true
						break
					}
				}
				if !isDisabled {
					notDisabled = append(notDisabled, cmdName)
				}
			}

			if len(notDisabled) > 0 {
				return &commands.CommandResponse{
					Content: fmt.Sprintf("These commands aren't disabled:\n\n%s\n\nHow tf do you plan to enable already enabled commands??",
						formatCommandList(notDisabled)),
				}, nil
			}

			// Enable commands
			dbBot, ok := ctx.Bot.(interface {
				GetDB() interface {
					EnableCommands(guildID string, commands []string) error
				}
			})

			if ok {
				if err := dbBot.GetDB().EnableCommands(ctx.Message.GuildID, normalizedArgs); err != nil {
					return &commands.CommandResponse{Content: "Failed to enable commands"}, nil
				}
			}

			return &commands.CommandResponse{
				Content: fmt.Sprintf("The following commands have been enabled successfully:\n\n%s",
					formatCommandList(normalizedArgs)),
			}, nil
		},
	})
}
