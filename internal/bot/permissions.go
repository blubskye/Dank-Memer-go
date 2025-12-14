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
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/dankmemer/bot/internal/utils"
)

// Permission GIF URLs for help messages
var permissionGifs = map[int64]string{
	discordgo.PermissionSendMessages:    "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionEmbedLinks:      "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionAttachFiles:     "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionReadMessages:    "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionManageMessages:  "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionVoiceConnect:    "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionVoiceSpeak:      "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionAddReactions:    "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionManageRoles:     "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionManageChannels:  "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionManageGuild:     "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionAdministrator:   "https://i.imgur.com/REDACTEd.gif",
	discordgo.PermissionUseExternalEmojis: "https://i.imgur.com/REDACTEd.gif",
}

// Permission names for display
var permissionNames = map[int64]string{
	discordgo.PermissionSendMessages:    "Send Messages",
	discordgo.PermissionEmbedLinks:      "Embed Links",
	discordgo.PermissionAttachFiles:     "Attach Files",
	discordgo.PermissionReadMessages:    "Read Messages",
	discordgo.PermissionManageMessages:  "Manage Messages",
	discordgo.PermissionVoiceConnect:    "Connect (Voice)",
	discordgo.PermissionVoiceSpeak:      "Speak (Voice)",
	discordgo.PermissionAddReactions:    "Add Reactions",
	discordgo.PermissionManageRoles:     "Manage Roles",
	discordgo.PermissionManageChannels:  "Manage Channels",
	discordgo.PermissionManageGuild:     "Manage Server",
	discordgo.PermissionAdministrator:   "Administrator",
	discordgo.PermissionUseExternalEmojis: "Use External Emojis",
}

func (b *Bot) checkPermissions(s *discordgo.Session, m *discordgo.MessageCreate, required []int64) bool {
	if len(required) == 0 {
		return true
	}

	// Get bot's permissions in this channel
	perms, err := s.State.UserChannelPermissions(s.State.User.ID, m.ChannelID)
	if err != nil {
		// Try to fetch directly
		perms, err = s.UserChannelPermissions(s.State.User.ID, m.ChannelID)
		if err != nil {
			b.Logger.Error().Err(err).Msg("Failed to get permissions")
			return true // Fail open, let Discord API reject if needed
		}
	}

	var missing []int64
	for _, perm := range required {
		if perms&perm == 0 {
			missing = append(missing, perm)
		}
	}

	if len(missing) == 0 {
		return true
	}

	// Send permission error
	b.sendPermissionError(s, m, missing)
	return false
}

func (b *Bot) sendPermissionError(s *discordgo.Session, m *discordgo.MessageCreate, missing []int64) {
	var permNames []string
	var gifURL string

	for _, perm := range missing {
		if name, ok := permissionNames[perm]; ok {
			permNames = append(permNames, name)
		}
		if gif, ok := permissionGifs[perm]; ok && gifURL == "" {
			gifURL = gif
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: "I'm missing permissions!",
		Description: fmt.Sprintf("I need the following permissions to run this command:\n`%s`\n\n"+
			"Please give me these permissions and try again.",
			strings.Join(permNames, "`, `")),
		Color: utils.RandomColor(),
	}

	if gifURL != "" {
		embed.Image = &discordgo.MessageEmbedImage{URL: gifURL}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// CheckUserPermissions checks if a user has specific permissions
func (b *Bot) CheckUserPermissions(s *discordgo.Session, guildID, userID string, required int64) bool {
	perms, err := s.State.UserChannelPermissions(userID, guildID)
	if err != nil {
		return false
	}
	return perms&required == required
}

// IsAdmin checks if a user has administrator permissions
func (b *Bot) IsAdmin(s *discordgo.Session, guildID, userID string) bool {
	return b.CheckUserPermissions(s, guildID, userID, discordgo.PermissionAdministrator)
}

// IsModerator checks if a user has moderation permissions
func (b *Bot) IsModerator(s *discordgo.Session, guildID, userID string) bool {
	modPerms := int64(discordgo.PermissionManageMessages | discordgo.PermissionKickMembers | discordgo.PermissionBanMembers)
	return b.CheckUserPermissions(s, guildID, userID, modPerms)
}
