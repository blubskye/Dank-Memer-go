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
			Triggers:    []string{"disable"},
			Description: "Use this command to disable commands you do not wish for your server to use",
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
					Content: "You are not authorized to use this command. You must have `Manage Server` to disable commands.",
				}, nil
			}

			if len(ctx.Args) == 0 {
				return &commands.CommandResponse{
					Content: fmt.Sprintf("Specify a command to disable, or multiple.\n\nExample: `%s disable meme trigger shitsound` or `%s disable meme`",
						ctx.GuildConfig.Prefix, ctx.GuildConfig.Prefix),
					Reply: true,
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
					cmdName = "nsfw" // Special case for NSFW category
				} else {
					continue // Invalid command
				}

				if !seen[cmdName] {
					seen[cmdName] = true
					normalizedArgs = append(normalizedArgs, cmdName)
				}
			}

			if len(normalizedArgs) == 0 {
				return &commands.CommandResponse{
					Content: "No valid commands specified.",
					Reply:   true,
				}, nil
			}

			// Check which are already disabled
			var alreadyDisabled []string
			for _, cmdName := range normalizedArgs {
				for _, disabled := range ctx.GuildConfig.DisabledCommands {
					if disabled == cmdName {
						alreadyDisabled = append(alreadyDisabled, cmdName)
						break
					}
				}
			}

			if len(alreadyDisabled) > 0 {
				return &commands.CommandResponse{
					Content: fmt.Sprintf("These commands are already disabled:\n\n%s\n\nHow tf do you plan to disable already disabled commands??",
						formatCommandList(alreadyDisabled)),
				}, nil
			}

			// Add to disabled commands
			dbBot, ok := ctx.Bot.(interface {
				GetDB() interface {
					DisableCommands(guildID string, commands []string) error
				}
			})

			if ok {
				if err := dbBot.GetDB().DisableCommands(ctx.Message.GuildID, normalizedArgs); err != nil {
					return &commands.CommandResponse{Content: "Failed to disable commands"}, nil
				}
			}

			return &commands.CommandResponse{
				Content: fmt.Sprintf("The following commands have been disabled successfully:\n\n%s",
					formatCommandList(normalizedArgs)),
			}, nil
		},
	})
}

func formatCommandList(cmds []string) string {
	var formatted []string
	for _, cmd := range cmds {
		formatted = append(formatted, fmt.Sprintf("`%s`", cmd))
	}
	return strings.Join(formatted, ", ")
}
