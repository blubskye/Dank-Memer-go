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

package voice

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Manager handles voice connections and audio playback
type Manager struct {
	session     *discordgo.Session
	connections map[string]*Connection
	mu          sync.RWMutex
}

// Connection represents an active voice connection
type Connection struct {
	GuildID    string
	ChannelID  string
	Voice      *discordgo.VoiceConnection
	Playing    bool
	StopChan   chan struct{}
}

// NewManager creates a new voice manager
func NewManager(session *discordgo.Session) *Manager {
	return &Manager{
		session:     session,
		connections: make(map[string]*Connection),
	}
}

// IsPlaying checks if audio is playing in a guild
func (m *Manager) IsPlaying(guildID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if conn, ok := m.connections[guildID]; ok {
		return conn.Playing
	}
	return false
}

// PlayAudio plays an audio file in a voice channel
func (m *Manager) PlayAudio(guildID, channelID, audioPath string) error {
	m.mu.Lock()

	// Check for existing connection
	if conn, ok := m.connections[guildID]; ok && conn.Playing {
		m.mu.Unlock()
		return ErrAlreadyPlaying
	}

	// Join voice channel
	vc, err := m.session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		m.mu.Unlock()
		return err
	}

	conn := &Connection{
		GuildID:   guildID,
		ChannelID: channelID,
		Voice:     vc,
		Playing:   true,
		StopChan:  make(chan struct{}),
	}
	m.connections[guildID] = conn
	m.mu.Unlock()

	// Play audio in goroutine
	go m.streamAudio(conn, audioPath)

	return nil
}

// Stop stops audio playback in a guild
func (m *Manager) Stop(guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, ok := m.connections[guildID]; ok {
		close(conn.StopChan)
		conn.Playing = false
	}
}

// Disconnect disconnects from a voice channel
func (m *Manager) Disconnect(guildID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, ok := m.connections[guildID]; ok {
		if conn.Voice != nil {
			conn.Voice.Disconnect()
		}
		delete(m.connections, guildID)
	}
	return nil
}

func (m *Manager) streamAudio(conn *Connection, audioPath string) {
	defer func() {
		m.mu.Lock()
		conn.Playing = false
		if conn.Voice != nil {
			conn.Voice.Disconnect()
		}
		delete(m.connections, conn.GuildID)
		m.mu.Unlock()
	}()

	// Open audio file
	file, err := os.Open(audioPath)
	if err != nil {
		return
	}
	defer file.Close()

	// Wait for voice connection to be ready
	time.Sleep(250 * time.Millisecond)

	// Start speaking
	conn.Voice.Speaking(true)
	defer conn.Voice.Speaking(false)

	// Stream audio
	// For now, we'll use a simple raw opus streaming approach
	// A more complete implementation would use DCA format
	m.streamOpus(conn, file)
}

func (m *Manager) streamOpus(conn *Connection, reader io.Reader) {
	// Read opus frames and send them
	// This is a simplified implementation
	// For production, use a proper DCA decoder like github.com/jonas747/dca

	buf := make([]byte, 960*2*2) // Opus frame buffer

	for {
		select {
		case <-conn.StopChan:
			return
		default:
		}

		n, err := reader.Read(buf)
		if err != nil {
			return
		}

		if n == 0 {
			return
		}

		select {
		case conn.Voice.OpusSend <- buf[:n]:
		case <-conn.StopChan:
			return
		case <-time.After(time.Second):
			return
		}
	}
}

// Error types
type VoiceError string

func (e VoiceError) Error() string { return string(e) }

const (
	ErrAlreadyPlaying VoiceError = "already playing audio in this guild"
	ErrNotConnected   VoiceError = "not connected to voice"
)
