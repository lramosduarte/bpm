package helpers

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func FollowUp(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: content,
	}); err != nil {
		slog.Warn("failed to send follow-up message", slog.String("error", err.Error()))
	}
}
