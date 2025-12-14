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
	"math/rand"
)

// Discord embed colors
var EmbedColors = []int{
	0x77dd77, // pastel green
	0x836953, // pastel brown
	0x89cff0, // baby blue
	0x99c5c4, // pastel teal
	0x9adedb, // baby turquoise
	0xaa9499, // pastel mauve
	0xaaf0d1, // mint
	0xb2fba5, // light green
	0xb39eb5, // pastel purple
	0xbdb0d0, // light purple
	0xbee7a5, // lime
	0xbefd73, // neon lime
	0xc1c6fc, // periwinkle
	0xc6a4a4, // dusty rose
	0xcb99c9, // orchid
	0xdea5a4, // pastel red
	0xf49ac2, // pink
	0xfea3aa, // light red
	0xfea3aa, // pastel salmon
	0xff6961, // coral
	0xffb7ce, // baby pink
	0xffdac1, // peach
	0xfdfd96, // pastel yellow
	0xffe4b5, // moccasin
}

func RandomColor() int {
	return EmbedColors[rand.Intn(len(EmbedColors))]
}
