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

// BaseCommand provides a simple implementation of the Command interface
type BaseCommand struct {
	Properties CommandProps
	Handler    func(*CommandContext) (*CommandResponse, error)
}

func NewBaseCommand(props CommandProps, handler func(*CommandContext) (*CommandResponse, error)) *BaseCommand {
	return &BaseCommand{
		Properties: props,
		Handler:    handler,
	}
}

func (c *BaseCommand) Props() CommandProps {
	props := c.Properties

	// Apply defaults
	if props.Cooldown == 0 {
		props.Cooldown = 3000
	}
	if props.Usage == "" {
		props.Usage = "{command}"
	}

	return props
}

func (c *BaseCommand) Run(ctx *CommandContext) (*CommandResponse, error) {
	// Check missing args
	if c.Properties.MissingArgs != "" && len(ctx.Args) == 0 {
		return &CommandResponse{Content: c.Properties.MissingArgs}, nil
	}

	return c.Handler(ctx)
}
