package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func PingHandler(name string, d discordServer) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
}
