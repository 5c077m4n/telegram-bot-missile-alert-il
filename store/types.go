package store

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChatID    int64              `bson:"chat_id"`
	City      string             `bson:"city"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type City struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Name            string             `bson:"name"`
	Area            string             `bson:"area"`
	SecondsToBunker int                `bson:"seonds_to_bunker"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}
