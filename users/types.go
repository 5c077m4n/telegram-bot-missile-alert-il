package users

import "time"

type User struct {
	ChatID    int64
	City      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
