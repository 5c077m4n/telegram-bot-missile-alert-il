package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"telegram-bot-missile-alert-il/alerts"
	"telegram-bot-missile-alert-il/cities"
	"telegram-bot-missile-alert-il/handler"
	"telegram-bot-missile-alert-il/history"

	"github.com/go-telegram/bot"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	token, found := os.LookupEnv("TELEGRAM_BOT_API_TOKEN")
	if !found {
		slog.Error("Telegram bot API token not found")
		os.Exit(1)
	}

	opts := []bot.Option{bot.WithDefaultHandler(handler.HandleMessage)}
	b, err := bot.New(token, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create bot: %v\n", err)
		os.Exit(1)
	}
	allCities, _ := cities.FetchCities()
	slog.Info(
		"Got all cities",
		"count", len(allCities),
	)
	allHistory, _ := history.FetchHistory()
	slog.Info(
		"Got entire history",
		"count", len(allHistory),
	)
	go alerts.PollAlerts(ctx, b)

	slog.Info("Telegram bot is starting...")
	b.Start(ctx)
}
