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

package utils

import (
	"encoding/json"
	"os"
)

type Assets struct {
	Kill           []string `json:"kill"`
	TrumpPhotos    []string `json:"trumpPhotos"`
	TrumpResponses []string `json:"trumpResponses"`
	Roast          []string `json:"roast"`
	Discord        []string `json:"discord"`
}

var LoadedAssets *Assets

func LoadAssets(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		// Initialize with empty data
		LoadedAssets = &Assets{
			Kill:           []string{"$author kills $mention"},
			TrumpPhotos:    []string{},
			TrumpResponses: []string{"huge", "the greatest", "fake news"},
			Roast:          []string{"You're a disappointment"},
			Discord:        []string{},
		}
		return nil
	}

	LoadedAssets = &Assets{}
	return json.Unmarshal(data, LoadedAssets)
}

func GetKillMessage() string {
	if LoadedAssets == nil || len(LoadedAssets.Kill) == 0 {
		return "$author kills $mention"
	}
	return RandomInArray(LoadedAssets.Kill)
}

func GetRoastMessage() string {
	if LoadedAssets == nil || len(LoadedAssets.Roast) == 0 {
		return "You're a disappointment"
	}
	return RandomInArray(LoadedAssets.Roast)
}

func GetTrumpResponse() string {
	if LoadedAssets == nil || len(LoadedAssets.TrumpResponses) == 0 {
		return "huge"
	}
	return RandomInArray(LoadedAssets.TrumpResponses)
}

func GetTrumpPhoto() string {
	if LoadedAssets == nil || len(LoadedAssets.TrumpPhotos) == 0 {
		return ""
	}
	return RandomInArray(LoadedAssets.TrumpPhotos)
}

func GetDiscordMeme() string {
	if LoadedAssets == nil || len(LoadedAssets.Discord) == 0 {
		return ""
	}
	return RandomInArray(LoadedAssets.Discord)
}
