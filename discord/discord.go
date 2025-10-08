package discord

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v11"
)

const (
	guildID = "" // TODO config this from interation
)

type Config struct {
	Token string `env:"DISCORD_TOKEN,required"`
}

type Discord struct {
	Config   *Config
	Session  *discordgo.Session
	commands []*Command
}

func New() *Discord {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("failed to parse env to config", slog.String("error", err.Error()))
	}
	return &Discord{Config: &cfg, commands: commands}
}

func (d *Discord) Close() {
	slog.Info("👋 Shutting down bot.")
	for _, cmd := range d.commands {
		slog.Debug("cleaning up command", slog.String("command", cmd.definition.Name))
		err := d.Session.ApplicationCommandDelete(d.Session.State.User.ID, guildID, cmd.registered.ID)
		if err != nil {
			slog.Error("Cannot delete", slog.String("command", cmd.definition.Name), slog.String("error", err.Error()))
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

	for _, cmd := range d.commands {
		discordCmd, err := d.Session.ApplicationCommandCreate(d.Session.State.User.ID, guildID, cmd.definition)
		if err != nil {
			return fmt.Errorf("cannot create '%s' command: %w", discordCmd.Name, err)
		}
		d.Session.AddHandler(cmd.handler(cmd.definition.Name))
		cmd.registered = discordCmd
	}

	slog.Info("🤖 Bot is now running. Press CTRL-C to exit.")

	return nil
}
