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
	"time"
)

func (db *Database) GetCooldown(command, userID string) (int64, error) {
	var expiresAt int64
	err := db.pool.QueryRow(`
		SELECT expires_at FROM cooldowns
		WHERE user_id = ? AND command = ?`, userID, command).
		Scan(&expiresAt)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	return expiresAt, err
}

func (db *Database) SetCooldown(command, userID string, durationMs int64) error {
	expiresAt := time.Now().UnixMilli() + durationMs
	_, err := db.pool.Exec(`
		INSERT INTO cooldowns (user_id, command, expires_at)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE expires_at = ?`,
		userID, command, expiresAt, expiresAt)
	return err
}

func (db *Database) ClearCooldown(command, userID string) error {
	_, err := db.pool.Exec(`
		DELETE FROM cooldowns WHERE user_id = ? AND command = ?`, userID, command)
	return err
}

func (db *Database) ClearAllCooldowns(userID string) error {
	_, err := db.pool.Exec(`DELETE FROM cooldowns WHERE user_id = ?`, userID)
	return err
}

func (db *Database) CleanupExpiredCooldowns() (int64, error) {
	result, err := db.pool.Exec(`
		DELETE FROM cooldowns WHERE expires_at < ?`, time.Now().UnixMilli())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// IsOnCooldown returns remaining time in ms if on cooldown, 0 if not
func (db *Database) IsOnCooldown(command, userID string) (int64, error) {
	expiresAt, err := db.GetCooldown(command, userID)
	if err != nil {
		return 0, err
	}
	if expiresAt == 0 {
		return 0, nil
	}

	now := time.Now().UnixMilli()
	if expiresAt > now {
		return expiresAt - now, nil
	}
	return 0, nil
}
