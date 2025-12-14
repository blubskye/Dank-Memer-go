# Dank Memer - Go Port

A Go port of the Dank Memer Discord bot, originally written in Node.js with Eris.

## Requirements

- Go 1.21+
- MariaDB/MySQL 8.0+
- Discord Bot Token

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/blubskye/Dank-Memer-go.git
cd Dank-Memer-go
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up the database

```sql
CREATE DATABASE dankmemer;
CREATE USER 'memer'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON dankmemer.* TO 'memer'@'localhost';
FLUSH PRIVILEGES;
```

Then run the migrations:

```bash
mysql -u root -p dankmemer < migrations/001_initial.sql
```

### 4. Configure the bot

Copy the example config and edit it:

```bash
cp config.yaml.example config.yaml
```

Edit `config.yaml` with your settings:

```yaml
token: "YOUR_DISCORD_BOT_TOKEN"
default_prefix: "pls"

database:
  host: "localhost"
  port: 3306
  user: "memer"
  password: "your_password"
  name: "dankmemer"
```

### 5. Build and run

```bash
go build -o memer ./cmd/memer
./memer
```

Or run directly:

```bash
go run ./cmd/memer
```

## Project Structure

```
├── cmd/memer/main.go          # Entry point
├── internal/
│   ├── bot/                   # Core bot logic
│   │   ├── bot.go             # Bot struct and lifecycle
│   │   ├── handler.go         # Message handler
│   │   └── permissions.go     # Permission checking
│   ├── commands/              # Command system
│   │   ├── types.go           # Command interfaces
│   │   ├── registry.go        # Command registration
│   │   ├── base.go            # BaseCommand
│   │   ├── image.go           # ImageCommand
│   │   ├── media.go           # MediaCommand
│   │   ├── reddit.go          # RedditCommand
│   │   ├── voice.go           # VoiceCommand
│   │   ├── animal/            # Animal commands (6)
│   │   ├── currency/          # Currency commands (2)
│   │   ├── fun/               # Fun commands (19)
│   │   ├── image/             # Image manipulation (23)
│   │   ├── meme/              # Meme commands (9)
│   │   ├── nsfw/              # NSFW commands (5)
│   │   ├── text/              # Text commands (2)
│   │   ├── utility/           # Utility commands (14)
│   │   └── voice/             # Voice commands (8)
│   ├── database/              # Database layer
│   ├── external/              # External API clients
│   ├── utils/                 # Utilities
│   └── voice/                 # Voice management
├── assets/                    # Static assets (audio, JSON data)
├── migrations/                # SQL migrations
└── config.yaml                # Configuration
```

## Commands (88 total)

### Text Commands (2)
- `clap` - Say something with clap emojis
- `emojify` - Convert text to emoji letters

### Fun Commands (19)
- `roast` - Roast someone
- `kill` - Kill someone (virtually)
- `mock` - Mock text in SpOnGeBoB case
- `asktrump` - Ask Trump a question
- `dankrate` - Rate how dank something is
- `waifu` - Rate how good of a waifu you are
- `repeat` - Repeat your message
- `chucknorris` - Chuck Norris jokes
- `xkcd` - Random XKCD comic
- `showerthoughts` - Shower thoughts from Reddit
- `tifu` - TIFU stories from Reddit
- `greentext` - Greentext stories
- `freenitro` - Free nitro (rickroll)
- `gay` - Rate how gay you are
- `google` - LMGTFY link
- `vent` - Send a message to the bot owner
- `comics` - Comics from Reddit
- `facepalm` - Facepalm images

### Currency Commands (2)
- `coins` - Check your coin balance
- `daily` - Collect daily coins

### Utility Commands (14)
- `help` - Show help
- `ping` - Ping the bot
- `prefix` - Change server prefix
- `stats` - Bot statistics
- `invite` - Bot invite link
- `patreon` - Patreon link
- `credits` - Bot credits
- `changes` - Changelog
- `website` - Bot website
- `enable` - Enable commands
- `disable` - Disable commands
- `clean` - Clean bot messages
- `dm` - DM a user (owner only)
- `source` - Get source code link (AGPL compliance)

### Animal Commands (6)
- `pupper` - Random dog picture
- `kitty` - Random cat picture
- `birb` - Random bird picture
- `aww` - Cute animals from Reddit
- `redpanda` - Random red panda
- `lizzyboi` - Random lizard

### Meme Commands (9)
- `meme` - Random meme from Reddit
- `joke` - Random joke from Reddit
- `discordmeme` - Discord memes
- `meirl` - Me IRL memes
- `pun` - Puns from Reddit
- `shitpost` - Shitposts from Reddit
- `wholesome` - Wholesome memes
- `prequel` - Prequel memes

### Image Manipulation (23)
- `magik` - Magik effect
- `ban` - Ban overlay
- `jail` - Jail overlay
- `trigger` - Triggered GIF
- `tweet` - Trump tweet generator
- `trash` - Trash overlay
- `spank` - Spank image
- `salty` - Salty GIF
- `search` - Google search image
- `shit` - Shit on something
- `rip` - RIP image
- `pride` - Pride flag overlay
- `hitler` - Worse than Hitler
- `failure` - Failure image
- `egg` - Egg birth image
- `disability` - Disability image
- `delete` - Delete this image
- `cry` - Cry image
- `cancer` - Cancer image
- `byemom` - Bye mom image
- `brazzers` - Brazzers logo
- `batslap` - Batman slap
- `b1nzy` - b1nzy cat meme

### NSFW Commands (5)
- `boobies` - NSFW content
- `booty` - NSFW content
- `4k` - 4K NSFW content
- `gayp` - Gay NSFW content
- `porngif` - NSFW GIFs

### Voice Commands (8)
- `airhorn` - Play airhorn
- `fart` - Play fart sound
- `boo` - Play scare sound
- `knock` - Play knock sound
- `mememusic` - Play meme sounds
- `mlg` - Play MLG sounds
- `oof` - Play Roblox oof
- `stop` - Stop audio playback

## Adding New Commands

Create a new file in the appropriate command category folder:

```go
// Dank Memer - A Discord bot
// Copyright (C) 2025 Dank Memer
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package fun

import (
    "github.com/dankmemer/bot/internal/bot"
    "github.com/dankmemer/bot/internal/commands"
)

func init() {
    bot.Register(&commands.BaseCommand{
        Properties: commands.CommandProps{
            Triggers:    []string{"example"},
            Description: "An example command",
            Category:    "Fun Commands",
        },
        Handler: func(ctx *commands.CommandContext) (*commands.CommandResponse, error) {
            return &commands.CommandResponse{Content: "Hello!"}, nil
        },
    })
}
```

## Authors

### Original Dank Memer (Node.js)
- **Melmsie** - *Initial work and owner* - [GitHub](https://github.com/melmsie)
- **Aetheryx** - *I like trains* - [GitHub](https://github.com/Aetheryx)
- **CyberRonin** - *Melmsie is my lover* - [GitHub](https://github.com/TheCyberRonin)
- **DaJuukes** - *Haha yes* - [GitHub](https://github.com/DaJuukes)
- **Kromatic** - *Mayonnaise is an instrument* - [GitHub](https://github.com/Devoxin)

### Go Port
- **blubskye** - *Go port* - [GitHub](https://github.com/blubskye)

## Acknowledgments

- **Stupid Cat** - *Original author of the backend code for the trigger command* - [GitHub](https://github.com/Ratismal)
- **Samoxive** - *Help with debugging stupid node errors and learning js* - [GitHub](https://github.com/Samoxive)
- **Kodehawa** - *Emotional support when dealing with stupid Discord issues and users* - [GitHub](https://github.com/Kodehawa)
- **Sporks** - *Went through 800 meme templates to give everyone a better experience*

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

This means:
- You can view, modify, and redistribute this source code
- If you run a modified version of this software as a network service, you must make the source code available to users
- Any modifications must be released under the same license
- You must give appropriate credit to the original authors

See the [LICENSE](LICENSE) file for the full license text, or visit https://www.gnu.org/licenses/agpl-3.0.html

### Source Code Command

This bot includes a `source` command (`pls source`) that provides a link to the source code, ensuring compliance with AGPL-3.0 requirements for network services.
