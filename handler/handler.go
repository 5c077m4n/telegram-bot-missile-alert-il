// Package handler - handles telegram input messages
package handler

import (
	"context"
	"log/slog"
	"slices"

	"telegram-bot-missile-alert-il/users"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	availableCities []string,
) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text
	slog.Info(
		"Recieved message",
		"chat ID", chatID,
		"message", text,
	)

	if !slices.Contains(availableCities, text) {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   text + "\nis a stupid city, please choose a better one",
		})
		if err != nil {
			slog.Error(
				"Error sending re-choose city message",
				"chat ID", chatID,
				"error", err,
			)
		}
		return
	}

	user := users.GetUser(chatID)

	if user.City == "" {
		user.City = text
		if err := users.UpdateUserCity(chatID, user.City); err != nil {
			slog.Error(
				"Error updating user city",
				"chat ID", chatID,
				"error", err,
			)
			return
		}

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "You are now subscribed to alerts for: " + user.City,
		})
		if err != nil {
			slog.Error(
				"Error sending subscription successful message",
				"chat ID", chatID,
				"error", err,
			)
		}
	} else {
		_, err := b.SendMessage(
			ctx,
			&bot.SendMessageParams{ChatID: chatID, Text: "Updated your city to: " + user.City},
		)
		if err != nil {
			slog.Error(
				"Error sending city update message",
				"chat ID", chatID,
				"error", err,
			)
		}

	}
}
