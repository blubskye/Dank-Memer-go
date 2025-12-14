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

package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dankmemer/bot/internal/bot"
	"github.com/dankmemer/bot/internal/utils"

	// Import command packages to register them
	_ "github.com/dankmemer/bot/internal/commands/animal"
	_ "github.com/dankmemer/bot/internal/commands/currency"
	_ "github.com/dankmemer/bot/internal/commands/fun"
	_ "github.com/dankmemer/bot/internal/commands/image"
	_ "github.com/dankmemer/bot/internal/commands/meme"
	_ "github.com/dankmemer/bot/internal/commands/nsfw"
	_ "github.com/dankmemer/bot/internal/commands/text"
	_ "github.com/dankmemer/bot/internal/commands/utility"
	_ "github.com/dankmemer/bot/internal/commands/voice"
)

func main() {
	// Parse flags
	configPath := flag.String("config", "config.yaml", "Path to config file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Setup logging
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Str("config", *configPath).Msg("Loading configuration")

	// Load config
	cfg, err := utils.LoadConfig(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	if cfg.Token == "" {
		log.Fatal().Msg("No bot token provided in config")
	}

	log.Info().Str("version", cfg.Version).Msg("Starting Dank Memer")

	// Load assets
	if err := utils.LoadAssets("assets/arrays.json"); err != nil {
		log.Warn().Err(err).Msg("Failed to load assets, using defaults")
	}

	// Create and start bot
	b, err := bot.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create bot")
	}

	// Register commands
	bot.RegisterCommands(b)

	if err := b.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

	log.Info().Int("commands", b.Commands.Count()).Msg("Bot started successfully")

	// Wait for shutdown
	b.Wait()

	// Cleanup
	if err := b.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Error during shutdown")
	}
}
