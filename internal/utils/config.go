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
	"github.com/spf13/viper"
)

type Config struct {
	Token         string   `mapstructure:"token"`
	DefaultPrefix string   `mapstructure:"default_prefix"`
	Version       string   `mapstructure:"version"`
	Devs          []string `mapstructure:"devs"`
	Premium       bool     `mapstructure:"premium"`
	PremiumGuilds []string `mapstructure:"premium_guilds"`
	VentChannel   string   `mapstructure:"vent_channel"`

	Database DatabaseConfig `mapstructure:"database"`
	Sharding ShardingConfig `mapstructure:"sharding"`
	APIs     APIsConfig     `mapstructure:"apis"`
	Webhooks WebhooksConfig `mapstructure:"webhooks"`
	Voice    VoiceConfig    `mapstructure:"voice"`
	URLs     URLsConfig     `mapstructure:"urls"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	MaxConns int    `mapstructure:"max_conns"`
}

type ShardingConfig struct {
	Enabled    bool `mapstructure:"enabled"`
	ShardCount int  `mapstructure:"shard_count"`
}

type APIsConfig struct {
	ImgenKey  string `mapstructure:"imgen_key"`
	ImgenURL  string `mapstructure:"imgen_url"`
	RedditURL string `mapstructure:"reddit_url"`
}

type WebhooksConfig struct {
	Shard   WebhookInfo `mapstructure:"shard"`
	Cluster WebhookInfo `mapstructure:"cluster"`
}

type WebhookInfo struct {
	ID    string `mapstructure:"id"`
	Token string `mapstructure:"token"`
}

type VoiceConfig struct {
	AudioPath string `mapstructure:"audio_path"`
}

type URLsConfig struct {
	Invite  string `mapstructure:"invite"`
	Support string `mapstructure:"support"`
	Patreon string `mapstructure:"patreon"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Environment variable overrides
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MEMER")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Set defaults
	if cfg.DefaultPrefix == "" {
		cfg.DefaultPrefix = "pls"
	}
	if cfg.Database.MaxConns == 0 {
		cfg.Database.MaxConns = 25
	}
	if cfg.Sharding.ShardCount == 0 {
		cfg.Sharding.ShardCount = 1
	}

	return &cfg, nil
}
