package telebot

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"telegram-bot-missile-alert-il/cities"
	"telegram-bot-missile-alert-il/users"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func handleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `*Welcome\!*

Type /city to choose your city 😎`,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		slog.Error("Could not send welcome message", "error", err)
	}
}

func handleCityChange(ctx context.Context, b *bot.Bot, update *models.Update, store *users.Store) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	newCity := strings.TrimLeft(update.Message.Text, "/city ")
	newCity = strings.Trim(newCity, " ")
	slog.Info("Recieved message", "chat ID", chatID, "new city", newCity)

	if newCity == "" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Can't be empty...",
		})
		if err != nil {
			slog.Error("Error sending non-empty warngin message", "chat ID", chatID, "error", err)
		}
		return
	}

	if availableCities, err := cities.FetchAllCityNames(); err == nil {
		slog.Info("Got all cities", "count", len(availableCities))

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
	} else {
		slog.Error("Could not fetch city list", "chat ID", chatID, "error", err)
	}

	user, err := store.GetUser(chatID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "My bad, couldn't get your user",
		})
		slog.Error("Could not get user", "chatID", chatID, "error", err)
		return
	}

	user, err = store.UpdateUserCity(user.ChatID, newCity)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "My bad, couldn't update your city",
		})
		slog.Error("Couldn't update user's city", "chatID", chatID)
		return
	}
	slog.Info("Updated user info", "user", user)

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{ChatID: chatID, Text: "Updated your city to: " + user.City},
	)
	if err != nil {
		slog.Error("Error sending city update message", "chat ID", chatID, "error", err)
	}
}
