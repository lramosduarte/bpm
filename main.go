package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"bpm/discord"
	"bpm/torrents/clients/qbitorrent"

	"github.com/caarlos0/env/v11"
)

type config struct {
	Discord discord.Config
	Torrent qbitorrent.Config
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("failed to parse env to config", slog.String("error", err.Error()))
		panic(err)
	}

	qbTorrent := qbitorrent.New(&cfg.Torrent)
	bot := discord.New(qbTorrent, &cfg.Discord)
	if err := bot.Start(); err != nil {
		panic(err)
	}
	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop
}
