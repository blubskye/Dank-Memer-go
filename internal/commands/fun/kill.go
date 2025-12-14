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

	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"kill", "murder"},
			Description: "Sick of someone? Easy! Just kill them!",
			Usage:       "{command} @user",
			Category:    "Fun Commands",
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Check for self-kill
			if len(ctx.Args) > 0 && strings.ToLower(ctx.Args[0]) == "me" {
				return &commands.CommandResponse{Content: "Ok you're dead. Please tag someone else to kill."}, nil
			}

			if len(ctx.Message.Mentions) == 0 {
				return &commands.CommandResponse{Content: "Ok you're dead. Please tag someone else to kill."}, nil
			}

			if ctx.Message.Mentions[0].ID == ctx.Message.Author.ID {
				return &commands.CommandResponse{Content: "Ok you're dead. Please tag someone else to kill."}, nil
			}

			msg := utils.GetKillMessage()

			// Replace placeholders
			msg = strings.ReplaceAll(msg, "$mention", ctx.Message.Mentions[0].Username)
			msg = strings.ReplaceAll(msg, "$author", ctx.Message.Author.Username)

			return &commands.CommandResponse{Content: msg}, nil
		},
	})
}
