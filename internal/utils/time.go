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
	"fmt"
	"strings"
	"time"
)

// FormatDuration formats a duration in milliseconds to a human-readable string
func FormatDuration(ms int64) string {
	d := time.Duration(ms) * time.Millisecond

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string

	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}

	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}

	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	if seconds > 0 || len(parts) == 0 {
		if seconds == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}

	if len(parts) == 1 {
		return parts[0]
	}

	// Join all but last with ", " and last with " and "
	var result strings.Builder
	for i, part := range parts {
		if i == 0 {
			result.WriteString(part)
		} else if i == len(parts)-1 {
			result.WriteString(" and ")
			result.WriteString(part)
		} else {
			result.WriteString(", ")
			result.WriteString(part)
		}
	}

	return result.String()
}

// FormatDurationShort formats a duration in milliseconds to a short string (e.g., "5m 30s")
func FormatDurationShort(ms int64) string {
	d := time.Duration(ms) * time.Millisecond

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// ParseTimeString parses a time string like "1h", "30m", "1d" and returns milliseconds
func ParseTimeString(s string) (int64, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return d.Milliseconds(), nil
}
