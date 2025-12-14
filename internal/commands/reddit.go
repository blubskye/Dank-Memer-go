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
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dankmemer/bot/internal/utils"
)

// RedditPostType defines what type of posts to fetch
type RedditPostType string

const (
	RedditPostTypeImage RedditPostType = "image"
	RedditPostTypeText  RedditPostType = "text"
)

// RedditCommand handles commands that fetch content from Reddit
type RedditCommand struct {
	Properties CommandProps
	Endpoint   string         // Reddit endpoint (e.g., "/r/memes/top.json")
	PostType   RedditPostType // Type of posts to fetch
}

func NewRedditCommand(props CommandProps, endpoint string, postType RedditPostType) *RedditCommand {
	return &RedditCommand{
		Properties: props,
		Endpoint:   endpoint,
		PostType:   postType,
	}
}

func (c *RedditCommand) Props() CommandProps {
	props := c.Properties

	// Apply defaults
	if props.Cooldown == 0 {
		props.Cooldown = 3000
	}

	// Add required permissions
	props.Permissions = append(props.Permissions, discordgo.PermissionEmbedLinks)

	return props
}

func (c *RedditCommand) Run(ctx *CommandContext) (*CommandResponse, error) {
	// Get the bot interface
	bot, ok := ctx.Bot.(redditBot)
	if !ok {
		return nil, fmt.Errorf("bot does not support Reddit fetching")
	}

	// Fetch posts
	var posts []redditPost
	var err error

	client := bot.GetRedditClient()
	if c.PostType == RedditPostTypeImage {
		posts, err = fetchImagePosts(client, c.Endpoint)
	} else {
		posts, err = fetchTextPosts(client, c.Endpoint)
	}

	if err != nil {
		return &CommandResponse{Content: fmt.Sprintf("Error fetching from Reddit: %s", err.Error())}, nil
	}

	if len(posts) == 0 {
		return &CommandResponse{Content: "No posts found!"}, nil
	}

	// Get index for this guild/command
	cmdName := c.Properties.Triggers[0]
	index := bot.GetRedditIndex(ctx.Message.GuildID, cmdName)

	if index >= len(posts) {
		index = 0
	}

	post := posts[index]
	bot.IncrementRedditIndex(ctx.Message.GuildID, cmdName, len(posts))

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title: utils.TruncateString(post.Title, 256),
		URL:   fmt.Sprintf("https://reddit.com%s", post.Permalink),
		Color: utils.RandomColor(),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("posted by %s", post.Author),
		},
	}

	if c.PostType == RedditPostTypeImage {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: post.URL,
		}
	} else {
		embed.Description = utils.TruncateString(post.Selftext, 2000)
	}

	return &CommandResponse{Embed: embed}, nil
}

// Interface to access bot's Reddit client
type redditBot interface {
	GetRedditClient() redditClient
	GetRedditIndex(guildID, command string) int
	IncrementRedditIndex(guildID, command string, maxIndex int) int
}

type redditClient interface {
	FetchPosts(endpoint string) ([]redditPost, error)
}

type redditPost struct {
	Title     string
	URL       string
	Permalink string
	Author    string
	Selftext  string
	PostHint  string
}

// Helper functions to convert from external package types
func fetchImagePosts(client redditClient, endpoint string) ([]redditPost, error) {
	posts, err := client.FetchPosts(endpoint)
	if err != nil {
		return nil, err
	}

	var imagePosts []redditPost
	for _, post := range posts {
		if post.PostHint == "image" || isRedditImageURL(post.URL) {
			imagePosts = append(imagePosts, post)
		}
	}
	return imagePosts, nil
}

func fetchTextPosts(client redditClient, endpoint string) ([]redditPost, error) {
	posts, err := client.FetchPosts(endpoint)
	if err != nil {
		return nil, err
	}

	var textPosts []redditPost
	for _, post := range posts {
		if post.Selftext != "" && len(post.Selftext) <= 2000 && len(post.Title) <= 256 {
			textPosts = append(textPosts, post)
		}
	}
	return textPosts, nil
}

func isRedditImageURL(url string) bool {
	return isImageURL(url)
}
