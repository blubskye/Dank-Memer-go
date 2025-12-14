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

type Donator struct {
	ID    string
	Level int
}

func (db *Database) IsDonator(userID string) (bool, error) {
	var exists int
	err := db.pool.QueryRow(`
		SELECT 1 FROM donators WHERE id = ?`, userID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *Database) GetDonator(userID string) (*Donator, error) {
	var d Donator
	err := db.pool.QueryRow(`
		SELECT id, level FROM donators WHERE id = ?`, userID).
		Scan(&d.ID, &d.Level)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (db *Database) GetDonatorLevel(userID string) (int, error) {
	d, err := db.GetDonator(userID)
	if err != nil {
		return 0, err
	}
	if d == nil {
		return 0, nil
	}
	return d.Level, nil
}

func (db *Database) SetDonator(userID string, level int) error {
	_, err := db.pool.Exec(`
		INSERT INTO donators (id, level) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE level = ?`, userID, level, level)
	return err
}

func (db *Database) RemoveDonator(userID string) error {
	_, err := db.pool.Exec(`DELETE FROM donators WHERE id = ?`, userID)
	return err
}
