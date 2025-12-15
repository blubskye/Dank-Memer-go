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
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/utils"
)

var startTime = time.Now()

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"stats", "info"},
			Description: "Returns information and statistics about Dank Memer.",
			Category:    "Utility Commands",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Get stats from database
			var guilds, users, channels int

			dbBot, ok := ctx.Bot.(interface {
				GetDB() interface {
					GetStats() (int, int, int, error)
				}
			})

			if ok {
				g, u, c, err := dbBot.GetDB().GetStats()
				if err == nil {
					guilds, users, channels = g, u, c
				}
			}

			// Fallback to session state if available
			if guilds == 0 {
				guilds = len(ctx.Session.State.Guilds)
				for _, g := range ctx.Session.State.Guilds {
					users += g.MemberCount
					channels += len(g.Channels)
				}
			}

			// Get command count
			cmdCount := 0
			cmdBot, ok := ctx.Bot.(interface {
				GetCommands() *commands.Registry
			})
			if ok {
				cmdCount = len(cmdBot.GetCommands().GetAll())
			}

			// Calculate uptime
			uptime := time.Since(startTime)

			// Get memory stats
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)

			avgSize := 0
			if guilds > 0 {
				avgSize = users / guilds
			}

			largeGuilds := 0
			for _, g := range ctx.Session.State.Guilds {
				if g.MemberCount >= 100 {
					largeGuilds++
				}
			}

			embed := &discordgo.MessageEmbed{
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "Server Statistics",
						Value: fmt.Sprintf("%s servers\n%d average server size\n%d large servers",
							formatNumber(guilds), avgSize, largeGuilds),
						Inline: true,
					},
					{
						Name: "Various Statistics",
						Value: fmt.Sprintf("%s uptime\n%s users\n%d commands currently",
							formatDuration(uptime), formatNumber(users), cmdCount),
						Inline: true,
					},
					{
						Name: "System Statistics",
						Value: fmt.Sprintf("%.1f MB memory\n%s platform\nGo %s",
							float64(mem.Alloc)/1024/1024, runtime.GOOS, runtime.Version()),
						Inline: true,
					},
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}

func formatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return strconv.Itoa(n)
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
