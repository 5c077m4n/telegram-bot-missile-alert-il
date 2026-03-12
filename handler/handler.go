// Package handler - handles telegram input messages
package handler

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"telegram-bot-missile-alert-il/users"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleCityChange(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	availableCities []string,
) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	newCity := strings.TrimLeft(update.Message.Text, "/city ")
	slog.Info("Recieved message", "chat ID", chatID, "new city", newCity)

	if !slices.Contains(availableCities, newCity) {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   newCity + "\nis a stupid city, please choose a better one",
		})
		if err != nil {
			slog.Error("Error sending re-choose city message", "chat ID", chatID, "error", err)
		}
		return
	}

	user := users.GetUser(chatID)
	user, err := users.UpdateUserCity(user.ChatID, newCity)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "My bad, couldn't update your city",
		})
		slog.Error("Could not find user", "chatID", user.ChatID)
		return
	}

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{ChatID: chatID, Text: "Updated your city to: " + user.City},
	)
	if err != nil {
		slog.Error("Error sending city update message", "chat ID", chatID, "error", err)
	}

}
