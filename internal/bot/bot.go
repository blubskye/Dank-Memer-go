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

package bot

import (
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/database"
	"github.com/dankmemer/bot/internal/external"
	"github.com/dankmemer/bot/internal/utils"
	"github.com/dankmemer/bot/internal/voice"
)

type Bot struct {
	Session  *discordgo.Session
	Config   *utils.Config
	DB       *database.Database
	Commands *commands.Registry
	Logger   zerolog.Logger

	// External clients
	ImageGen     *external.ImageGenClient
	RedditClient *external.RedditClient
	VoiceManager *voice.Manager

	// Runtime state
	MentionRegex  *regexp.Regexp
	RedditIndexes map[string]map[string]int // guildID -> command -> index
	redditMu      sync.RWMutex

	// Shutdown handling
	shutdownChan chan struct{}
}

func New(cfg *utils.Config) (*Bot, error) {
	// Initialize logger
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Str("component", "bot").Logger()

	// Create Discord session
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	// Create database connection
	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Session:       session,
		Config:        cfg,
		DB:            db,
		Commands:      commands.NewRegistry(),
		Logger:        logger,
		RedditIndexes: make(map[string]map[string]int),
		shutdownChan:  make(chan struct{}),
	}

	// Initialize external clients
	bot.ImageGen = external.NewImageGenClient(cfg.APIs.ImgenURL, cfg.APIs.ImgenKey)
	bot.RedditClient = external.NewRedditClient(cfg.APIs.RedditURL)
	bot.VoiceManager = voice.NewManager(session)

	return bot, nil
}

func (b *Bot) Start() error {
	b.Logger.Info().Msg("Starting bot...")

	// Set intents
	b.Session.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsMessageContent

	// Register event handlers
	b.Session.AddHandler(b.handleReady)
	b.Session.AddHandler(b.handleMessageCreate)
	b.Session.AddHandler(b.handleGuildCreate)
	b.Session.AddHandler(b.handleGuildDelete)

	// Open connection
	if err := b.Session.Open(); err != nil {
		return err
	}

	b.Logger.Info().Msg("Bot connected to Discord")

	return nil
}

func (b *Bot) Wait() {
	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case <-sc:
		b.Logger.Info().Msg("Received shutdown signal")
	case <-b.shutdownChan:
		b.Logger.Info().Msg("Shutdown requested")
	}
}

func (b *Bot) Shutdown() error {
	b.Logger.Info().Msg("Shutting down bot...")

	close(b.shutdownChan)

	if err := b.Session.Close(); err != nil {
		b.Logger.Error().Err(err).Msg("Error closing Discord session")
	}

	if err := b.DB.Close(); err != nil {
		b.Logger.Error().Err(err).Msg("Error closing database")
	}

	b.Logger.Info().Msg("Bot shutdown complete")
	return nil
}

func (b *Bot) handleReady(s *discordgo.Session, r *discordgo.Ready) {
	b.Logger.Info().
		Str("username", r.User.Username).
		Str("id", r.User.ID).
		Int("guilds", len(r.Guilds)).
		Msg("Bot is ready")

	// Compile mention regex
	b.MentionRegex = regexp.MustCompile(`^<@!?` + r.User.ID + `>\s*`)

	// Set status
	s.UpdateGameStatus(0, b.Config.DefaultPrefix+" help | "+b.Config.Version)

	// Update stats
	b.updateStats()
}

func (b *Bot) handleGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	b.Logger.Debug().Str("guild", g.ID).Str("name", g.Name).Msg("Joined guild")

	// Create guild config if not exists
	_, err := b.DB.GetOrCreateGuild(g.ID, b.Config.DefaultPrefix)
	if err != nil {
		b.Logger.Error().Err(err).Str("guild", g.ID).Msg("Failed to create guild config")
	}

	b.updateStats()
}

func (b *Bot) handleGuildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	b.Logger.Debug().Str("guild", g.ID).Msg("Left guild")

	// Optionally delete guild config
	// b.DB.DeleteGuild(g.ID)

	b.updateStats()
}

func (b *Bot) updateStats() {
	guilds := len(b.Session.State.Guilds)
	var users, channels int

	for _, g := range b.Session.State.Guilds {
		users += g.MemberCount
		channels += len(g.Channels)
	}

	stats := &database.BotStats{
		Guilds:   guilds,
		Users:    users,
		Channels: channels,
		Shards:   1,
	}

	if err := b.DB.UpdateStats(stats); err != nil {
		b.Logger.Error().Err(err).Msg("Failed to update stats")
	}
}

// GetRedditIndex returns the current index for a Reddit command in a guild
func (b *Bot) GetRedditIndex(guildID, command string) int {
	b.redditMu.RLock()
	defer b.redditMu.RUnlock()

	if guildIndexes, ok := b.RedditIndexes[guildID]; ok {
		if idx, ok := guildIndexes[command]; ok {
			return idx
		}
	}
	return 0
}

// IncrementRedditIndex increments and returns the new index for a Reddit command
func (b *Bot) IncrementRedditIndex(guildID, command string, maxIndex int) int {
	b.redditMu.Lock()
	defer b.redditMu.Unlock()

	if _, ok := b.RedditIndexes[guildID]; !ok {
		b.RedditIndexes[guildID] = make(map[string]int)
	}

	idx := b.RedditIndexes[guildID][command]
	idx++
	if idx >= maxIndex {
		idx = 0
	}
	b.RedditIndexes[guildID][command] = idx

	return idx
}
