package history

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"telegram-bot-missile-alert-il/store"
	"time"

	"github.com/go-telegram/bot"
)

func fetchAlerts() ([]*Alert, error) {
	req, err := http.NewRequest(
		"GET",
		"https://www.oref.org.il/WarningMessages/alert/History/AlertsHistory.json",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Could not close request body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code from Pikud HaOref History API: %d",
			resp.StatusCode,
		)
	}

	var alerts []*Alert
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

func notifyUsers(ctx context.Context, b *bot.Bot, s *store.Store, alerts []*Alert) error {
	allUsers, err := s.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	var errs []error
	for _, user := range allUsers {
		if user.City == "" {
			continue
		}

		for _, alert := range alerts {
			if alert.ShouldSend(user.City) {
				slog.Debug("Found a match", "user", user, "history item", alert)

				_, err := b.SendMessage(
					ctx,
					&bot.SendMessageParams{ChatID: user.ChatID, Text: alert.Title},
				)
				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func Poll(ctx context.Context, b *bot.Bot, s *store.Store) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			alerts, err := fetchAlerts()
			if err != nil {
				continue
			}

			slog.Debug("Recieved historical alerts", "count", len(alerts))
			if err := notifyUsers(ctx, b, s, alerts); err != nil {
				slog.Error("could not notify user", "error", err, "alerts", alerts)
			}

		}
	}
}
