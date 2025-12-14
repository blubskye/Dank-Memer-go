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
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/dankmemer/bot/internal/commands"
	"github.com/dankmemer/bot/internal/database"
	"github.com/dankmemer/bot/internal/utils"
)

func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bots
	if m.Author.Bot {
		return
	}

	// Ignore DMs
	if m.GuildID == "" {
		return
	}

	// Check if blocked
	blocked, _ := b.DB.IsUserOrGuildBlocked(m.Author.ID, m.GuildID)
	if blocked {
		return
	}

	// Premium check
	if b.Config.Premium && !utils.Contains(b.Config.PremiumGuilds, m.GuildID) {
		if strings.HasPrefix(strings.ToLower(m.Content), b.Config.DefaultPrefix) {
			s.ChannelMessageSend(m.ChannelID,
				"This server is not a premium activated server. Want it activated? https://patreon.com/dank")
		}
		return
	}

	// Get guild config
	guildConfig, err := b.DB.GetOrCreateGuild(m.GuildID, b.Config.DefaultPrefix)
	if err != nil {
		b.Logger.Error().Err(err).Str("guild", m.GuildID).Msg("Failed to get guild config")
		guildConfig = &defaultGuildConfig
	}

	// Parse prefix and content
	prefix, content, ok := b.parsePrefix(m.Content, guildConfig.Prefix)
	if !ok {
		// Check for greeting with mention
		if b.MentionRegex != nil && b.MentionRegex.MatchString(m.Content) {
			if strings.Contains(strings.ToLower(m.Content), "hello") {
				s.ChannelMessageSend(m.ChannelID,
					fmt.Sprintf("Hello, %s. My prefix is `%s`. Example: `%s meme`",
						m.Author.Username, guildConfig.Prefix, guildConfig.Prefix))
			}
		}
		return
	}

	// Parse command and arguments
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	cmdName := strings.ToLower(parts[0])
	args := parts[1:]

	// Find command
	cmd := b.Commands.Find(cmdName)
	if cmd == nil {
		return
	}

	props := cmd.Props()

	// Check owner-only
	if props.OwnerOnly && !utils.Contains(b.Config.Devs, m.Author.ID) {
		return
	}

	// Check if command is disabled
	if utils.Contains(guildConfig.DisabledCommands, props.Triggers[0]) {
		return
	}

	// Check if NSFW category is disabled
	if props.IsNSFW && utils.Contains(guildConfig.DisabledCommands, "nsfw") {
		return
	}

	// Check cooldown
	cooldown := props.Cooldown
	if cooldown == 0 {
		cooldown = 3000 // Default 3 seconds
	}

	remainingCD, _ := b.DB.IsOnCooldown(props.Triggers[0], m.Author.ID)
	if remainingCD > 0 {
		msg := props.CooldownMessage
		if msg == "" {
			msg = "stop spamming my commands dude, you have to wait {cooldown}"
		}
		msg = strings.Replace(msg, "{cooldown}", utils.FormatDuration(remainingCD), 1)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	// Check permissions
	if !b.checkPermissions(s, m, props.Permissions) {
		return
	}

	// Check NSFW channel
	if props.IsNSFW {
		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			channel, err = s.Channel(m.ChannelID)
		}
		if err != nil || !channel.NSFW {
			embed := &discordgo.MessageEmbed{
				Title:       "NSFW not allowed here",
				Description: "Use NSFW commands in a NSFW marked channel",
				Color:       utils.RandomColor(),
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}
	}

	// Create context
	ctx := &commands.CommandContext{
		Session:     s,
		Message:     m,
		Args:        args,
		CleanArgs:   b.resolveCleanArgs(s, m.Message, args),
		GuildConfig: guildConfig,
		Bot:         b,
	}

	// Execute command in goroutine
	go b.executeCommand(ctx, cmd, prefix)
}

func (b *Bot) parsePrefix(content, guildPrefix string) (prefix, rest string, ok bool) {
	lower := strings.ToLower(content)

	// Check mention prefix
	if b.MentionRegex != nil && b.MentionRegex.MatchString(content) {
		prefix = b.MentionRegex.FindString(content)
		rest = strings.TrimPrefix(content, prefix)
		rest = strings.TrimSpace(rest)
		return prefix, rest, true
	}

	// Check guild prefix
	if strings.HasPrefix(lower, strings.ToLower(guildPrefix)) {
		prefix = content[:len(guildPrefix)]
		rest = strings.TrimSpace(content[len(guildPrefix):])
		return prefix, rest, true
	}

	return "", "", false
}

func (b *Bot) executeCommand(ctx *commands.CommandContext, cmd commands.Command, prefix string) {
	defer func() {
		if r := recover(); r != nil {
			b.Logger.Error().
				Interface("panic", r).
				Str("command", cmd.Props().Triggers[0]).
				Str("user", ctx.Message.Author.ID).
				Str("guild", ctx.Message.GuildID).
				Msg("Command panicked")

			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID,
				"Something went wrong while executing that command. Please try again later.")
		}
	}()

	props := cmd.Props()

	// Log command execution
	b.Logger.Debug().
		Str("command", props.Triggers[0]).
		Str("user", ctx.Message.Author.ID).
		Str("guild", ctx.Message.GuildID).
		Strs("args", ctx.Args).
		Msg("Executing command")

	start := time.Now()

	// Execute the command
	resp, err := cmd.Run(ctx)

	duration := time.Since(start)

	if err != nil {
		b.Logger.Error().
			Err(err).
			Str("command", props.Triggers[0]).
			Strs("args", ctx.Args).
			Dur("duration", duration).
			Msg("Command error")

		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID,
			fmt.Sprintf("Something went wrong: `%s`", err.Error()))
		return
	}

	// Set cooldown after successful execution
	cooldown := props.Cooldown
	if cooldown == 0 {
		cooldown = 3000
	}
	b.DB.SetCooldown(props.Triggers[0], ctx.Message.Author.ID, cooldown)

	if resp == nil {
		return
	}

	// Send response
	b.sendResponse(ctx, resp)

	b.Logger.Debug().
		Str("command", props.Triggers[0]).
		Dur("duration", duration).
		Msg("Command completed")
}

func (b *Bot) sendResponse(ctx *commands.CommandContext, resp *commands.CommandResponse) {
	// Build message
	msg := &discordgo.MessageSend{}

	if resp.Content != "" {
		msg.Content = resp.Content
	}

	if resp.Embed != nil {
		// Set default color if not set
		if resp.Embed.Color == 0 {
			resp.Embed.Color = utils.RandomColor()
		}
		msg.Embeds = []*discordgo.MessageEmbed{resp.Embed}
	}

	if resp.File != nil {
		msg.Files = []*discordgo.File{resp.File}
	}

	if len(resp.Files) > 0 {
		msg.Files = resp.Files
	}

	if resp.Reply {
		msg.Reference = &discordgo.MessageReference{
			MessageID: ctx.Message.ID,
			ChannelID: ctx.Message.ChannelID,
			GuildID:   ctx.Message.GuildID,
		}
	}

	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, msg)
	if err != nil {
		b.Logger.Error().Err(err).Msg("Failed to send response")
	}
}

func (b *Bot) resolveCleanArgs(s *discordgo.Session, m *discordgo.Message, args []string) []string {
	clean := make([]string, len(args))
	for i, arg := range args {
		// Check if it's a mention
		if strings.HasPrefix(arg, "<@") && strings.HasSuffix(arg, ">") {
			// Extract user ID
			id := strings.TrimPrefix(arg, "<@")
			id = strings.TrimPrefix(id, "!")
			id = strings.TrimSuffix(id, ">")

			// Try to get username
			member, err := s.GuildMember(m.GuildID, id)
			if err == nil && member.User != nil {
				clean[i] = member.User.Username
				continue
			}
		}
		clean[i] = arg
	}
	return clean
}

var defaultGuildConfig = database.GuildConfig{
	Prefix:           "pls",
	DisabledCommands: []string{},
}
