package telebot

import (
	"context"
	"log/slog"
	"os"
	"telegram-bot-missile-alert-il/alerts"

	"github.com/go-telegram/bot"
)

func Setup(ctx context.Context) error {
	token, found := os.LookupEnv("TELEGRAM_BOT_API_TOKEN")
	if !found {
		return ErrTokenNotFound
	}

	b, err := bot.New(token)
	if err != nil {
		return err
	}
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		handleStart,
	)
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/city",
		bot.MatchTypePrefix,
		handleCityChange,
	)

	go alerts.PollAlerts(ctx, b)

	slog.Info("Telegram bot is starting...")
	b.Start(ctx)

	return nil
}
