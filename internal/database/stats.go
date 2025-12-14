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

package database

type BotStats struct {
	Guilds   int
	Users    int
	Channels int
	Shards   int
}

func (db *Database) GetStats() (*BotStats, error) {
	var stats BotStats
	err := db.pool.QueryRow(`
		SELECT guilds, users, channels, shards FROM stats WHERE id = 1`).
		Scan(&stats.Guilds, &stats.Users, &stats.Channels, &stats.Shards)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (db *Database) UpdateStats(stats *BotStats) error {
	_, err := db.pool.Exec(`
		INSERT INTO stats (id, guilds, users, channels, shards) VALUES (1, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE guilds = ?, users = ?, channels = ?, shards = ?`,
		stats.Guilds, stats.Users, stats.Channels, stats.Shards,
		stats.Guilds, stats.Users, stats.Channels, stats.Shards)
	return err
}
