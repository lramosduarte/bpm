# bpm

## What is this?

`bpm` is a self-hosted Discord bot designed to help you manage your PC remotely. Currently, it focuses on torrent management via qBittorrent, but the goal is to expand its capabilities to allow you to control and automate various aspects of your computer directly from Discord.


## Features Available

- Discord bot integration
- Manage torrents via Discord commands
   - Add torrents (currently supports qBittorrent)
- Environment-based configuration
- Graceful shutdown handling

### Planned Features
- PC management commands
- More torrent clients
- Custom automation via Discord

## How to Run

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd bpm
   ```
2. **Set up environment variables:**
   The bot uses environment variables for configuration. You can set them directly in your shell or use a `.env` file. See the configuration section below for details.
3. **Build and run:**
   ```bash
   go build -o bpm
   ./bpm
   ```

## How to Configure Discord

You need a Discord bot token and the bot must be invited to your server.

### Steps:
1. **Create a Discord bot:**
   - Go to [Discord Developer Portal](https://discord.com/developers/applications)
   - Create a new application and add a bot
   - Copy the bot token
2. **Invite the bot to your server:**
   - In the Developer Portal, go to "OAuth2 > URL Generator"
   - Select "bot" scope and set permissions as needed (e.g., Send Messages, Read Message History)
   - Use the generated URL to invite the bot
3. **Set environment variables for Discord config:**
   The following environment variables are required:
   - `DISCORD_TOKEN`: Your bot token

4. **Run slash command `/setup` in your discord server**
    This command will configure essential information, such as the guild ID, for your server.

Example `.env` file:
```env
DISCORD_TOKEN=your-bot-token
TORRENT_HOST=localhost
TORRENT_PORT=8080
TORRENT_CREDENTIALS_USERNAME=admin
TORRENT_CREDENTIALS_PASSWORD=adminadmin
```

## Configuration

The bot uses the following environment variables (see `main.go`):
- Discord config: `DISCORD_*`
- qBittorrent config: `TORRENT_HOST`, `TORRENT_PORT`, `TORRENT_CREDENTIALS_USERNAME`, `TORRENT_CREDENTIALS_PASSWORD`


## License

MIT


# Roadmap

- [-] Torrent management
- [ ] PC management commands
- [ ] GuildID from envvars
- [ ] Generate Docker Image
