package history

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"telegram-bot-missile-alert-il/poller"
)

func FetchAlerts() ([]*Alert, error) {
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

func Stream(ctx context.Context) (<-chan []*Alert, <-chan error) {
	p := poller.New(FetchAlerts, 2*time.Second)
	return p.Stream(ctx)
}
