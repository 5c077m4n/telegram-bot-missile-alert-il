package alerts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"telegram-bot-missile-alert-il/users"

	"github.com/go-telegram/bot"
)

var (
	lastAlertID  int64
	lastAlertMtx sync.Mutex
)

func isDuplicateAlert(id int64) bool {
	lastAlertMtx.Lock()
	defer lastAlertMtx.Unlock()

	if id == lastAlertID {
		return true
	}
	lastAlertID = id
	return false
}

func formatAlertMessage(alert *WarningAlert) string {
	messageParts := []string{
		"🚨 ALERT: " + alert.Title + " 🚨",
		"Cities: " + strings.Join(alert.Cities, ", "),
	}
	return strings.Join(messageParts, "\n\n")
}

func notifyUsers(ctx context.Context, b *bot.Bot, alert *WarningAlert) error {
	allUsers := users.GetAllUsers()

	var errs []error
	for _, user := range allUsers {
		if user.City == "" {
			continue
		}

		if users.ContainsCityArray(alert.Cities, user.City) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: user.ChatID,
				Text:   formatAlertMessage(alert),
			})
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func fetchAlert() (*WarningAlert, error) {
	req, err := http.NewRequest(
		"GET",
		"https://www.oref.org.il/WarningMessages/alert/alerts.json",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", "https://www.oref.org.il/11226-he/pakar.aspx")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
	)
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
		return nil, fmt.Errorf("unexpected status code from Pikud HaOref API %d", resp.StatusCode)
	}

	var alert *WarningAlert
	if err := json.NewDecoder(resp.Body).Decode(alert); err != nil {
		return nil, err
	}
	slog.Info("Got an alert", "value", alert)

	return alert, nil
}

func PollAlerts(ctx context.Context, b *bot.Bot) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			alert, err := fetchAlert()
			if err != nil {
				continue
			}

			slog.Info(
				"Recieved alert",
				"value", alert,
			)
			if !isDuplicateAlert(alert.ID) {
				if err := notifyUsers(ctx, b, alert); err != nil {
					slog.Error(
						"could not notify user",
						"error", err,
						"alert", alert,
					)
				}
			}
		}
	}
}
