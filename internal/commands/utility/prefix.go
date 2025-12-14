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
			Triggers:    []string{"prefix"},
			Description: "Change Dank Memer's prefix!",
			Usage:       "{command} <prefix of your choice>",
			Category:    "Utility",
			Cooldown:    5000,
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			b, ok := ctx.Bot.(prefixBot)
			if !ok {
				return nil, fmt.Errorf("cannot access bot")
			}

			// Check permissions
			perms, err := ctx.Session.State.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err == nil {
				hasManageGuild := perms&discordgo.PermissionManageServer != 0
				isDev := utils.Contains(b.GetConfig().Devs, ctx.Message.Author.ID)

				if !hasManageGuild && !isDev {
					return &commands.CommandResponse{
						Content: "You are not authorized to use this command. You must have `Manage Server` to change the prefix.",
					}, nil
				}
			}

			currentPrefix := ctx.GuildConfig.Prefix

			// No args - show current prefix
			if len(ctx.Args) == 0 {
				return &commands.CommandResponse{
					Content: fmt.Sprintf("What do you want your new prefix to be?\n\nExample: `%s prefix pepe`", currentPrefix),
				}, nil
			}

			newPrefix := strings.ToLower(strings.Join(ctx.Args, " "))

			// Check length
			if len(newPrefix) > 32 {
				return &commands.CommandResponse{
					Content: fmt.Sprintf("Your prefix can't be over 32 characters long. You're %d characters over the limit.", len(newPrefix)-32),
				}, nil
			}

			// Check if same
			if newPrefix == currentPrefix {
				return &commands.CommandResponse{
					Content: fmt.Sprintf("`%s` is already your current prefix.", currentPrefix),
				}, nil
			}

			// Update prefix
			err = b.GetDB().UpdateGuildPrefix(ctx.Message.GuildID, newPrefix)
			if err != nil {
				return nil, err
			}

			embed := &discordgo.MessageEmbed{
				Description: fmt.Sprintf("Prefix successfully changed to `%s`.", newPrefix),
				Color:       utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}

type prefixBot interface {
	GetDB() prefixDB
	GetConfig() *utils.Config
}

type prefixDB interface {
	UpdateGuildPrefix(guildID, prefix string) error
}
