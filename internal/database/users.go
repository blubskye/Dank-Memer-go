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

import (
	"database/sql"
)

func (db *Database) GetCoins(userID string) (int64, error) {
	var coins int64
	err := db.pool.QueryRow(`SELECT coins FROM users WHERE id = ?`, userID).Scan(&coins)
	if err == sql.ErrNoRows {
		// Create user with 0 coins
		_, err = db.pool.Exec(`INSERT INTO users (id, coins) VALUES (?, 0)`, userID)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return coins, nil
}

func (db *Database) AddCoins(userID string, amount int64) error {
	_, err := db.pool.Exec(`
		INSERT INTO users (id, coins) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE coins = coins + ?`, userID, amount, amount)
	return err
}

func (db *Database) RemoveCoins(userID string, amount int64) error {
	_, err := db.pool.Exec(`
		UPDATE users SET coins = GREATEST(0, CAST(coins AS SIGNED) - ?) WHERE id = ?`, amount, userID)
	return err
}

func (db *Database) SetCoins(userID string, amount int64) error {
	_, err := db.pool.Exec(`
		INSERT INTO users (id, coins) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE coins = ?`, userID, amount, amount)
	return err
}
