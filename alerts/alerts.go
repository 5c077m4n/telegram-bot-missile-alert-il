package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"telegram-bot-missile-alert-il/poller"
)

func FetchAlert() (*Alert, error) {
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

	var alert *Alert
	if err := json.NewDecoder(resp.Body).Decode(alert); err != nil {
		return nil, err
	}

	return alert, nil
}

func Stream(ctx context.Context) (<-chan *Alert, <-chan error) {
	p := poller.New(FetchAlert, 2*time.Second)
	return p.Stream(ctx)
}
