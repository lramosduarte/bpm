package commands

import (
	"bpm/discord/helpers"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func AddTorrentHandler(name string, d discordServer) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

		if i.ApplicationCommandData().Resolved == nil || len(i.ApplicationCommandData().Resolved.Attachments) == 0 {
			helpers.FollowUp(s, i, "❌ No file was uploaded.")
			return
		}

		fileOpt := i.ApplicationCommandData().Options[0].Value.(string)
		attachment := i.ApplicationCommandData().Resolved.Attachments[fileOpt]

		if attachment.ContentType != "application/x-bittorrent" {
			helpers.FollowUp(s, i, "❌ The provided file is not a valid torrent file.")
			return
		}

		if err := d.Torrent().AddTorrent(attachment.URL); err != nil {
			helpers.FollowUp(s, i, "❌ Failed to add the torrent to the download queue.")
			slog.Error("failed to add torrent", slog.String("error", err.Error()))
			return
		}

		slog.Debug("attachment received", slog.String("filename", attachment.Filename), slog.String("url", attachment.URL), slog.Int("size", attachment.Size))
		helpers.FollowUp(s, i, "Torrent added to the download queue! 🎉")
	})
}
