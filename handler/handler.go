// Package handler - handles telegram input messages
package handler

import (
	"context"
	"log/slog"

	"telegram-bot-missile-alert-il/users"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	user := users.GetUser(chatID)

	if user.City == "" {
		user.City = text
		if err := users.UpdateUserCity(chatID, user.City); err != nil {
			slog.Error(
				"Error updating user city for chat",
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
				"Error sending message to chat",
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
				"Error sending message to chat",
				"chat ID", chatID,
				"error", err,
			)
		}

	}
}
