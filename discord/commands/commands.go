package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

const (
	CommandSetup = "setup"
)

type ClientTorrent interface {
	AddTorrent(url string) error
}

type discordServer interface {
	IsConfigured() bool
	Setup(*discordgo.InteractionCreate) error
	Torrent() ClientTorrent
}

type SlashHandler func(name string, discord discordServer) func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Command struct {
	Definition *discordgo.ApplicationCommand
	Registered *discordgo.ApplicationCommand
	Handler    SlashHandler
}

var Commands = []*Command{
	{
		Definition: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "check the bot's responsiveness",
		},
		Handler: PingHandler,
	},
	{
		Definition: &discordgo.ApplicationCommand{
			Name:        CommandSetup,
			Description: "Setup the bot for this server",
		},
		Handler: SetupHandler,
	},
	{
		Definition: &discordgo.ApplicationCommand{
			Name:        "add-torrent",
			Description: "Add a torrent to the download queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "from-file",
					Description: "Add a torrent from a file",
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Required:    false,
				},
				// {
				// 	Name:        "from-url",
				// 	Description: "Add a torrent from a url",
				// 	Type:        discordgo.ApplicationCommandOptionString,
				// 	Required:    false,
				// },
			},
		},
		Handler: AddTorrentHandler,
	},
}

func WithCommandNameCheck(name string, d discordServer, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name != name {
			return
		}
		if name != CommandSetup && !d.IsConfigured() {
			slog.Warn("bot is not configured for this server", slog.String("command", name), slog.String("guildID", i.GuildID))
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ Bot is not configured for this server. Please run the setup command.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		handler(s, i)
	}
}
