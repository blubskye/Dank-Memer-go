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

package nsfw

import (
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
)

func init() {
	bot.Register(&commands.MediaCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"booty", "ass", "butt"},
			Description: "See some thicc booties!",
			Category:    "NSFW Commands",
			IsNSFW:      true,
		},
		RequestURL: "https://boob.bot/api/v2/img/ass",
		JSONKey:    "url",
		Title:      "Here, take some booty.",
		Message:    "Free nudes thanks to boobbot & tom <3",
		TokenKey:   "porn",
	})
}
