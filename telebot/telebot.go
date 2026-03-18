package telebot

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"strings"

	"telegram-bot-missile-alert-il/store"

	"github.com/5c077m4n/pikud-haoref-api-go/alerts"
	"github.com/5c077m4n/pikud-haoref-api-go/history"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func formatAlertMessage(alert *alerts.Alert) string {
	messageParts := []string{
		"🚨 ALERT: " + alert.Title + " 🚨",
		"Cities: " + strings.Join(alert.Cities, ", "),
	}
	return strings.Join(messageParts, "\n\n")
}

func notifyAlertUsers(ctx context.Context, b *bot.Bot, s *store.Store, alert *alerts.Alert) error {
	if !alert.ShouldSend() {
		return nil
	}

	allUsers, err := s.GetAllUsers(ctx)
	if err != nil {
		slog.Error("Could not get all users", "error", err)
		return err
	}

	for _, user := range allUsers {
		if user.City == "" {
			continue
		}

		if slices.Contains(alert.Cities, user.City) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: user.ChatID,
				Text:   formatAlertMessage(alert),
			})
			if err != nil {
				slog.Error("Could not send alert message", "error", err, "chatID", user.ChatID)
			}
		}
	}

	return nil
}

func notifyHistoryUsers(
	ctx context.Context,
	b *bot.Bot,
	s *store.Store,
	historyAlerts []*history.Alert,
) error {
	allUsers, err := s.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range allUsers {
		if user.City == "" {
			continue
		}

		for _, alert := range historyAlerts {
			if alert.ShouldSend(user.City) {
				slog.Debug("Found a match", "user", user, "history item", alert)

				_, err := b.SendMessage(
					ctx,
					&bot.SendMessageParams{ChatID: user.ChatID, Text: alert.Title},
				)
				if err != nil {
					slog.Error(
						"Could not send history message",
						"error",
						err,
						"chatID",
						user.ChatID,
					)
				}
			}
		}
	}

	return nil
}

func Setup(ctx context.Context, s *store.Store) error {
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
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/where",
		bot.MatchTypeExact,
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			handleWhere(ctx, bot, update, s)
		},
	)

	go func() {
		alertResultsCh, alertErrorsCh := alerts.Stream(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-alertErrorsCh:
				continue
			case alert := <-alertResultsCh:
				slog.Debug("Received alert", "value", alert)
				if err := notifyAlertUsers(ctx, b, s, alert); err != nil {
					slog.Error("could not notify user", "error", err, "alert", alert)
				}
			}
		}
	}()

	go func() {
		historyResultsCh, historyErrorsCh := history.Stream(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-historyErrorsCh:
				continue
			case alerts := <-historyResultsCh:
				slog.Debug("Received historical alerts", "count", len(alerts))
				if err := notifyHistoryUsers(ctx, b, s, alerts); err != nil {
					slog.Error("could not notify user", "error", err, "alerts", alerts)
				}
			}
		}
	}()

	slog.Info("Telegram bot is starting...")
	b.Start(ctx)

	return nil
}
