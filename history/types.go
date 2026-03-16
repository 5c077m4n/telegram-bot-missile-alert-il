package history

import (
	"log/slog"
	"time"
)

type HistoryItemDate struct{ time.Time }

func (hDate *HistoryItemDate) UnmarshalJSON(bytes []byte) error {
	parsedDate, err := time.Parse("\"2006-01-02 15:04:05\"", string(bytes))
	if err != nil {
		slog.Error(
			"Could not parse date",
			"original date", string(bytes),
			"error", err,
		)
		return err
	}

	hDate.Time = parsedDate
	return nil
}

type HistoryItem struct {
	Date     HistoryItemDate `json:"alertDate"`
	Title    string          `json:"title"`
	City     string          `json:"data"`
	Category int             `json:"category"`
}
