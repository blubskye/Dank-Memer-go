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
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/external"
	"github.com/dankmemer/bot/internal/utils"
)

type xkcdComic struct {
	Num   int    `json:"num"`
	Title string `json:"title"`
	Alt   string `json:"alt"`
	Img   string `json:"img"`
	Link  string `json:"link"`
}

func init() {
	bot.Register(&commands.BaseCommand{
		Properties: commands.CommandProps{
			Triggers:    []string{"xkcd"},
			Description: "Grabs a random comic from xkcd",
			Category:    "Fun Commands",
			Permissions: []int64{discordgo.PermissionEmbedLinks},
		},
		Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
			// Get latest comic number
			latestData, err := external.Get("https://xkcd.com/info.0.json", nil)
			if err != nil {
				return &commands.CommandResponse{Content: "Couldn't fetch xkcd"}, nil
			}

			var latest xkcdComic
			if err := json.Unmarshal(latestData, &latest); err != nil {
				return &commands.CommandResponse{Content: "Couldn't parse xkcd data"}, nil
			}

			// Get random comic
			randomNum := rand.Intn(latest.Num) + 1
			comicData, err := external.Get(fmt.Sprintf("https://xkcd.com/%d/info.0.json", randomNum), nil)
			if err != nil {
				return &commands.CommandResponse{Content: "Couldn't fetch comic"}, nil
			}

			var comic xkcdComic
			if err := json.Unmarshal(comicData, &comic); err != nil {
				return &commands.CommandResponse{Content: "Couldn't parse comic data"}, nil
			}

			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					Name: fmt.Sprintf("Comic #%d %s", comic.Num, comic.Title),
					URL:  fmt.Sprintf("https://xkcd.com/%d", comic.Num),
				},
				Description: comic.Alt,
				Image: &discordgo.MessageEmbedImage{
					URL: comic.Img,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "https://xkcd.com",
				},
				Color: utils.RandomColor(),
			}

			return &commands.CommandResponse{Embed: embed}, nil
		},
	})
}
