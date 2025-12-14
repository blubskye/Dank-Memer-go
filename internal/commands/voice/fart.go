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

package voice

import (
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
)

func init() {
	bot.Register(&commands.VoiceCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"fart"},
			Description: "Fart in your voice channel!",
			Category:    "Voice Commands",
			Cooldown:    10000,
		},
		Directory:    "farts",
		FileCount:    5,
		Extension:    "opus",
		Reaction:     "💨",
		ExistingConn: "I can only fart one at a time! Please wait.",
	})
}
