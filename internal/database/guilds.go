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
	"encoding/json"
)

type GuildConfig struct {
	ID               string
	Prefix           string
	DisabledCommands []string
	Premium          bool
}

func (db *Database) GetGuild(guildID string) (*GuildConfig, error) {
	var cfg GuildConfig
	var disabledJSON []byte

	err := db.pool.QueryRow(`
		SELECT id, prefix, disabled_commands, premium
		FROM guilds WHERE id = ?`, guildID).
		Scan(&cfg.ID, &cfg.Prefix, &disabledJSON, &cfg.Premium)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(disabledJSON) > 0 {
		json.Unmarshal(disabledJSON, &cfg.DisabledCommands)
	}
	if cfg.DisabledCommands == nil {
		cfg.DisabledCommands = []string{}
	}

	return &cfg, nil
}

func (db *Database) CreateGuild(guildID, prefix string) (*GuildConfig, error) {
	_, err := db.pool.Exec(`
		INSERT INTO guilds (id, prefix, disabled_commands)
		VALUES (?, ?, '[]')
		ON DUPLICATE KEY UPDATE id=id`, guildID, prefix)
	if err != nil {
		return nil, err
	}
	return db.GetGuild(guildID)
}

func (db *Database) UpdateGuildPrefix(guildID, prefix string) error {
	_, err := db.pool.Exec(`
		UPDATE guilds SET prefix = ? WHERE id = ?`, prefix, guildID)
	return err
}

func (db *Database) UpdateGuildDisabledCommands(guildID string, disabled []string) error {
	disabledJSON, err := json.Marshal(disabled)
	if err != nil {
		return err
	}
	_, err = db.pool.Exec(`
		UPDATE guilds SET disabled_commands = ? WHERE id = ?`, disabledJSON, guildID)
	return err
}

func (db *Database) UpdateGuildPremium(guildID string, premium bool) error {
	_, err := db.pool.Exec(`
		UPDATE guilds SET premium = ? WHERE id = ?`, premium, guildID)
	return err
}

func (db *Database) DeleteGuild(guildID string) error {
	_, err := db.pool.Exec(`DELETE FROM guilds WHERE id = ?`, guildID)
	return err
}

func (db *Database) GetOrCreateGuild(guildID, defaultPrefix string) (*GuildConfig, error) {
	cfg, err := db.GetGuild(guildID)
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		return cfg, nil
	}
	return db.CreateGuild(guildID, defaultPrefix)
}

func (db *Database) DisableCommands(guildID string, commands []string) error {
	cfg, err := db.GetGuild(guildID)
	if err != nil {
		return err
	}
	if cfg == nil {
		return nil
	}

	// Add new commands to disabled list
	for _, cmd := range commands {
		found := false
		for _, existing := range cfg.DisabledCommands {
			if existing == cmd {
				found = true
				break
			}
		}
		if !found {
			cfg.DisabledCommands = append(cfg.DisabledCommands, cmd)
		}
	}

	return db.UpdateGuildDisabledCommands(guildID, cfg.DisabledCommands)
}

func (db *Database) EnableCommands(guildID string, commands []string) error {
	cfg, err := db.GetGuild(guildID)
	if err != nil {
		return err
	}
	if cfg == nil {
		return nil
	}

	// Remove commands from disabled list
	newDisabled := make([]string, 0, len(cfg.DisabledCommands))
	for _, existing := range cfg.DisabledCommands {
		shouldKeep := true
		for _, cmd := range commands {
			if existing == cmd {
				shouldKeep = false
				break
			}
		}
		if shouldKeep {
			newDisabled = append(newDisabled, existing)
		}
	}

	return db.UpdateGuildDisabledCommands(guildID, newDisabled)
}
