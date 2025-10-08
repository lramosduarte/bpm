package discord

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type SlashHandler func(name string) func(s *discordgo.Session, i *discordgo.InteractionCreate)

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
		handler: func(name string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			return WithCommandNameCheck(name, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			Name:        "add-torrent",
			Description: "Add a torrent to the download queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "from-file",
					Description: "Add a torrent from a file",
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Required:    false,
				},
				{
					Name:        "from-url",
					Description: "Add a torrent from a url",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
		handler: func(name string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			return WithCommandNameCheck(name, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
					s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
						Content: "❌ No file was uploaded.",
					})
					return
				}

				fileOpt := i.ApplicationCommandData().Options[0].Value.(string)
				attachment := i.ApplicationCommandData().Resolved.Attachments[fileOpt]

				if attachment.ContentType != "application/x-bittorrent" {
					s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
						Content: "❌ The provided file is not a valid torrent file.",
					})
					return
				}

				// Download the .torrent file from Discord
				torrentData, err := downloadFile(attachment.URL)
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
						Content: "❌ Failed to download the torrent file.",
						Flags:   discordgo.MessageFlagsEphemeral,
					})
					return
				}

				// Send it to qBittorrent
				if err = addTorrentToQBittorrent(torrentData); err != nil {
					s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
						Content: fmt.Sprintf("❌ Failed to add torrent to qBittorrent: %v", err),
						Flags:   discordgo.MessageFlagsEphemeral,
					})
					return
				}

				slog.Debug("attachment received", slog.String("filename", attachment.Filename), slog.String("url", attachment.URL), slog.Int("size", attachment.Size))
				s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
					Content: "Torrent added to the download queue! 🎉",
				})
			})
		},
	},
}

func WithCommandNameCheck(name string, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name != name {
			return
		}
		handler(s, i)
	}
}
