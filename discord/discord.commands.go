package discord

import (
	"bpm/discord/helpers"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

const (
	CommandSetup = "setup"
)

type SlashHandler func(name string, discord *Discord) func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Command struct {
	definition *discordgo.ApplicationCommand
	registered *discordgo.ApplicationCommand
	handler    SlashHandler
}

var commands = []*Command{
	{
		definition: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "check the bot's responsiveness",
		},
		handler: func(name string, d *Discord) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			return WithCommandNameCheck(name, d, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Title:   "🏓 Pong!",
						Content: "bot is alive and responding!",
					},
				}); err != nil {
					slog.Error("failed to respond to interaction", slog.String("error", err.Error()))
				}
			})
		},
	},
	{
		definition: &discordgo.ApplicationCommand{
			Name:        CommandSetup,
			Description: "Setup the bot for this server",
		},
		handler: func(name string, d *Discord) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			return WithCommandNameCheck(name, d, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if d.GuildID == "" {
					slog.Info("bot is being configured for the first time", slog.String("guildID", i.GuildID))
					d.GuildID = i.GuildID
				}
				if d.GuildID != i.GuildID {
					if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "❌ Bot is already configured and using a different guildID. Please contact the bot administrator.",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					}); err != nil {
						slog.Error("failed to respond to interaction", slog.String("error", err.Error()))
					}
					return
				}

				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "✅ Bot has been configured for this server!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}); err != nil {
					slog.Error("failed to respond to interaction", slog.String("error", err.Error()))
				}
			})
		},
	},
	{
		definition: &discordgo.ApplicationCommand{
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
		handler: func(name string, d *Discord) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			return WithCommandNameCheck(name, d, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Title:   "Receiving file...",
						Content: "Downloading and checking the file...",
					},
				}); err != nil {
					slog.Error("failed to respond to interaction", slog.String("error", err.Error()))
				}

				if len(i.ApplicationCommandData().Resolved.Attachments) == 0 {
					helpers.FollowUp(s, i, "❌ No file was uploaded.")
					return
				}

				fileOpt := i.ApplicationCommandData().Options[0].Value.(string)
				attachment := i.ApplicationCommandData().Resolved.Attachments[fileOpt]

				if attachment.ContentType != "application/x-bittorrent" {
					helpers.FollowUp(s, i, "❌ The provided file is not a valid torrent file.")
					return
				}

				if err := d.clientTorrent.AddTorrent(attachment.URL); err != nil {
					helpers.FollowUp(s, i, "❌ Failed to add the torrent to the download queue.")
					slog.Error("failed to add torrent", slog.String("error", err.Error()))
					return
				}

				slog.Debug("attachment received", slog.String("filename", attachment.Filename), slog.String("url", attachment.URL), slog.Int("size", attachment.Size))
				helpers.FollowUp(s, i, "Torrent added to the download queue! 🎉")
			})
		},
	},
}

func WithCommandNameCheck(name string, d *Discord, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name != name {
			return
		}
		if name != CommandSetup && d.GuildID == "" {
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
