package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"bpm/discord"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	bot := discord.New()
	if err := bot.Start(); err != nil {
		panic(err)
	}
	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop
}
