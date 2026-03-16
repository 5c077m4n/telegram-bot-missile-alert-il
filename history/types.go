package history

import (
	"log/slog"
	"sync"
	"time"
)

type AlertDate struct{ time.Time }

func (d *AlertDate) UnmarshalJSON(bytes []byte) error {
	parsedDate, err := time.Parse("\"2006-01-02 15:04:05\"", string(bytes))
	if err != nil {
		slog.Error("Could not parse date", "original date", string(bytes), "error", err)
		return err
	}

	d.Time = parsedDate
	return nil
}

type Alert struct {
	Date     AlertDate `json:"alertDate"`
	Title    string    `json:"title"`
	City     string    `json:"data"`
	Category int       `json:"category"`
}

var (
	lastSentAt   = time.Now()
	lastSentAtMu sync.Mutex
)

func (a *Alert) ShouldSend(city string) bool {
	if a.City != city {
		return false
	}

	lastSentAtMu.Lock()
	defer lastSentAtMu.Unlock()

	if a.Date.After(lastSentAt) {
		lastSentAt = a.Date.Time
		return true
	}

	return false
}
