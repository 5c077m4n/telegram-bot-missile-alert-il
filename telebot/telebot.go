package telebot

import (
	"context"
	"log/slog"
	"os"
	"telegram-bot-missile-alert-il/alerts"
	"telegram-bot-missile-alert-il/users"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Setup(ctx context.Context, s *users.Store) error {
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
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			handleCityChange(ctx, bot, update, s)
		},
	)

	go alerts.PollAlerts(ctx, b, s)

	slog.Info("Telegram bot is starting...")
	b.Start(ctx)

	return nil
}
