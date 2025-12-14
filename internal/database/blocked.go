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

type BlockType string

const (
	BlockTypeUser  BlockType = "user"
	BlockTypeGuild BlockType = "guild"
)

type BlockedEntry struct {
	ID     string
	Type   BlockType
	Reason string
}

func (db *Database) IsBlocked(id string) (bool, error) {
	var exists int
	err := db.pool.QueryRow(`
		SELECT 1 FROM blocked WHERE id = ?`, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *Database) IsUserOrGuildBlocked(userID, guildID string) (bool, error) {
	var exists int
	err := db.pool.QueryRow(`
		SELECT 1 FROM blocked WHERE id IN (?, ?)`, userID, guildID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *Database) Block(id string, blockType BlockType, reason string) error {
	_, err := db.pool.Exec(`
		INSERT INTO blocked (id, type, reason) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE reason = ?`, id, blockType, reason, reason)
	return err
}

func (db *Database) Unblock(id string) error {
	_, err := db.pool.Exec(`DELETE FROM blocked WHERE id = ?`, id)
	return err
}

func (db *Database) GetBlocked(id string) (*BlockedEntry, error) {
	var entry BlockedEntry
	err := db.pool.QueryRow(`
		SELECT id, type, reason FROM blocked WHERE id = ?`, id).
		Scan(&entry.ID, &entry.Type, &entry.Reason)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
