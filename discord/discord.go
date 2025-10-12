package discord

import (
	"bpm/discord/commands"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token string `env:"BPM_DISCORD_TOKEN,required"`
}

type Discord struct {
	Config        *Config
	Session       *discordgo.Session
	GuildID       string
	commands      []*commands.Command
	clientTorrent commands.ClientTorrent
}

func New(clientTorrent commands.ClientTorrent, config *Config) *Discord {
	return &Discord{Config: config, commands: commands.Commands, clientTorrent: clientTorrent}
}

func (d *Discord) Close() {
	slog.Info("👋 Shutting down bot.")
	for _, cmd := range d.commands {
		slog.Debug("cleaning up command", slog.String("command", cmd.Definition.Name))
		err := d.Session.ApplicationCommandDelete(d.Session.State.User.ID, d.GuildID, cmd.Registered.ID)
		if err != nil {
			slog.Error("Cannot delete", slog.String("command", cmd.Definition.Name), slog.String("error", err.Error()))
		}
	}
	if d.Session != nil {
		if err := d.Session.Close(); err != nil {
			slog.Error("failed to close discord session", slog.String("error", err.Error()))
		}
	}
}

func (d *Discord) Start() error {
	var err error
	d.Session, err = discordgo.New("Bot " + d.Config.Token)
	if err != nil {
		return fmt.Errorf("Error(%w) creating Discord session", err)
	}

	if err = d.Session.Open(); err != nil {
		return fmt.Errorf("Error(%w) opening connection", err)
	}

	if d.GuildID == "" {
		slog.Warn("Bot is not configured.")
	}

	for _, cmd := range d.commands {
		discordCmd, err := d.Session.ApplicationCommandCreate(d.Session.State.User.ID, d.GuildID, cmd.Definition)
		if err != nil {
			return fmt.Errorf("cannot create '%s' command: %w", discordCmd.Name, err)
		}
		d.Session.AddHandler(cmd.Handler(cmd.Definition.Name, d))
		cmd.Registered = discordCmd
	}

	slog.Info("🤖 Bot is now running. Press CTRL-C to exit.")

	return nil
}

func (d *Discord) IsConfigured() bool {
	return d.GuildID != ""
}

func (d *Discord) Setup(i *discordgo.InteractionCreate) error {
	if d.IsConfigured() {
		return fmt.Errorf("bot is already configured")
	}
	d.GuildID = i.GuildID
	slog.Info("bot has been configured", slog.String("guildID", d.GuildID))
	return nil
}

func (d *Discord) Torrent() commands.ClientTorrent {
	return d.clientTorrent
}
