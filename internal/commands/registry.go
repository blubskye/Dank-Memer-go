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

package commands

import (
	"strings"
	"sync"
)

// Registry holds all registered commands
type Registry struct {
	commands []Command
	mu       sync.RWMutex
}

// NewRegistry creates a new command registry
func NewRegistry() *Registry {
	return &Registry{
		commands: make([]Command, 0),
	}
}

// Register adds a command to the registry
func (r *Registry) Register(cmd Command) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands = append(r.commands, cmd)
}

// Find looks up a command by trigger name
func (r *Registry) Find(trigger string) Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	trigger = strings.ToLower(trigger)
	for _, cmd := range r.commands {
		for _, t := range cmd.Props().Triggers {
			if strings.ToLower(t) == trigger {
				return cmd
			}
		}
	}
	return nil
}

// GetAll returns all registered commands
func (r *Registry) GetAll() []Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Command, len(r.commands))
	copy(result, r.commands)
	return result
}

// GetByCategory returns all commands in a specific category
func (r *Registry) GetByCategory(category string) []Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Command
	category = strings.ToLower(category)
	for _, cmd := range r.commands {
		if strings.ToLower(cmd.Props().Category) == category {
			result = append(result, cmd)
		}
	}
	return result
}

// GetCategories returns a list of all unique categories
func (r *Registry) GetCategories() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categoryMap := make(map[string]bool)
	for _, cmd := range r.commands {
		cat := cmd.Props().Category
		if cat != "" {
			categoryMap[cat] = true
		}
	}

	result := make([]string, 0, len(categoryMap))
	for cat := range categoryMap {
		result = append(result, cat)
	}
	return result
}

// Count returns the number of registered commands
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.commands)
}
