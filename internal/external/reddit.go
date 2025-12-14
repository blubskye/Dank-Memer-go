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

package external

import (
	"fmt"
	"strings"
)

// RedditClient handles requests to Reddit's JSON API
type RedditClient struct {
	baseURL string
}

func NewRedditClient(baseURL string) *RedditClient {
	if baseURL == "" {
		baseURL = "https://www.reddit.com"
	}
	return &RedditClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

// RedditPost represents a single Reddit post
type RedditPost struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Permalink string `json:"permalink"`
	Author    string `json:"author"`
	Subreddit string `json:"subreddit"`
	Score     int    `json:"score"`
	NSFW      bool   `json:"over_18"`
	IsVideo   bool   `json:"is_video"`
	Selftext  string `json:"selftext"`
	PostHint  string `json:"post_hint"`
}

// RedditListing represents the JSON structure from Reddit
type RedditListing struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// FetchPosts fetches posts from a Reddit endpoint
// endpoint should be like "/r/memes/top.json" or "/u/kerdaloo/m/dankmemer/top/.json"
func (c *RedditClient) FetchPosts(endpoint string) ([]RedditPost, error) {
	url := c.baseURL + endpoint

	// Add default params if not present
	if !strings.Contains(url, "?") {
		url += "?limit=100"
	}

	var listing RedditListing
	err := GetJSON(url, map[string]string{
		"User-Agent": "DankMemer/1.0",
	}, &listing)
	if err != nil {
		return nil, err
	}

	posts := make([]RedditPost, 0, len(listing.Data.Children))
	for _, child := range listing.Data.Children {
		posts = append(posts, child.Data)
	}

	return posts, nil
}

// FetchImagePosts fetches only image posts from a Reddit endpoint
func (c *RedditClient) FetchImagePosts(endpoint string) ([]RedditPost, error) {
	posts, err := c.FetchPosts(endpoint)
	if err != nil {
		return nil, err
	}

	imagePosts := make([]RedditPost, 0)
	for _, post := range posts {
		if isImageURL(post.URL) || post.PostHint == "image" {
			imagePosts = append(imagePosts, post)
		}
	}

	return imagePosts, nil
}

// FetchTextPosts fetches only text/self posts from a Reddit endpoint
func (c *RedditClient) FetchTextPosts(endpoint string) ([]RedditPost, error) {
	posts, err := c.FetchPosts(endpoint)
	if err != nil {
		return nil, err
	}

	textPosts := make([]RedditPost, 0)
	for _, post := range posts {
		if post.Selftext != "" {
			textPosts = append(textPosts, post)
		}
	}

	return textPosts, nil
}

// GetFullURL returns the full Reddit URL for a post
func (p *RedditPost) GetFullURL() string {
	return fmt.Sprintf("https://reddit.com%s", p.Permalink)
}

func isImageURL(url string) bool {
	lower := strings.ToLower(url)
	return strings.HasSuffix(lower, ".jpg") ||
		strings.HasSuffix(lower, ".jpeg") ||
		strings.HasSuffix(lower, ".png") ||
		strings.HasSuffix(lower, ".gif") ||
		strings.Contains(lower, "i.redd.it") ||
		strings.Contains(lower, "i.imgur.com")
}
