package telebot

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"telegram-bot-missile-alert-il/cities"
	"telegram-bot-missile-alert-il/store"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func handleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      `*Welcome\!* Type /city to choose your city 😎`,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		slog.Error("Could not send welcome message", "error", err)
	}
}

func handleCityChange(ctx context.Context, b *bot.Bot, update *models.Update, s *store.Store) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	newCity := strings.TrimLeft(update.Message.Text, "/city ")
	newCity = strings.Trim(newCity, " ")
	slog.Debug("Recieved message", "chat ID", chatID, "new city", newCity)

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

	user, err := s.GetOrCreateUser(chatID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "My bad, couldn't get your user",
		})
		slog.Error("Could not get user", "chatID", chatID, "error", err)
		return
	}

	if user.City == newCity {
		msg := "You're already here..."
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: msg}); err != nil {
			slog.Error(
				"Error sending already here message",
				"chat ID", chatID,
				"error", err,
				"message", msg,
			)
		}
		return
	}

	if availableCities, err := cities.FetchAllCityNamesOnce(); err != nil {
		slog.Error("Could not fetch city list", "chat ID", chatID, "error", err)
	} else {
		if !slices.Contains(availableCities, newCity) {
			msg := strings.Builder{}
			msg.WriteString(newCity)
			msg.WriteString("\nis a stupid city, please choose a better one")

			if similars := fuzzy.Find(newCity, availableCities); len(similars) > 0 {
				msg.WriteString("\nMaybe you meant:\n")
				msg.WriteString(strings.Join(similars, "\n"))
			}

			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: msg.String()}); err != nil {
				slog.Error("Error sending re-choose city message", "chat ID", chatID, "error", err, "message", msg)
			}
			return
		}
	}

	user, err = s.UpdateUserCity(user.ChatID, newCity)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "My bad, couldn't update your city",
		})
		slog.Error("Couldn't update user's city", "chatID", chatID)
		return
	}
	slog.Debug("Updated user info", "user", user)

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{ChatID: chatID, Text: "Updated your city to: " + user.City},
	)
	if err != nil {
		slog.Error("Error sending city update message", "chat ID", chatID, "error", err)
	}
}

func handleWhere(ctx context.Context, b *bot.Bot, update *models.Update, s *store.Store) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	user, err := s.GetUser(chatID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "My bad, couldn't get your user. Did you choose a city yet?",
		})
		slog.Error("Could not get user", "chatID", chatID, "error", err)
		return
	}

	if user.City != "" {
		_, errSendMessage := b.SendMessage(
			ctx,
			&bot.SendMessageParams{ChatID: chatID, Text: "You live in " + user.City},
		)
		if errSendMessage != nil {
			slog.Error("Error sending current city", "chat ID", chatID, "error", err)
		}
	} else {
		_, errSendMessage := b.SendMessage(
			ctx,
			&bot.SendMessageParams{ChatID: chatID, Text: "Don't know 🤷"},
		)
		if errSendMessage != nil {
			slog.Error("Error sending empty empty city message", "chat ID", chatID, "error", err)
		}
	}
}
