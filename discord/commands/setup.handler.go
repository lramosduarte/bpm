package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func SetupHandler(name string, d discordServer) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return WithCommandNameCheck(name, d, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if !d.IsConfigured() {
			slog.Info("bot is being configured for the first time", slog.String("guildID", i.GuildID))
			if err := d.Setup(i); err != nil {
				slog.Error("failed to setup bot", slog.String("error", err.Error()))
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "❌ Failed to setup the bot. Please contact the bot administrator.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}); err != nil {
					slog.Error("failed to respond to interaction", slog.String("error", err.Error()))
				}
				return
			}
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
}
