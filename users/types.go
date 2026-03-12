package users

import "time"

type User struct {
	ChatID    int64     `bson:"chat_id"`
	City      string    `bson:"city"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
