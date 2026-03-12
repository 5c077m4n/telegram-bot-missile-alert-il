package users

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client          *mongo.Client
	usersCollection *mongo.Collection
}

func NewStore(uri, dbName string) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &Store{client: client, usersCollection: client.Database(dbName).Collection("users")}, nil
}

func (s *Store) GetUser(chatID int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := s.usersCollection.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		user := &User{ChatID: chatID, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		_, err := s.usersCollection.InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) UpdateUserCity(chatID int64, city string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"chat_id": chatID}
	update := bson.M{
		"$set": bson.M{"city": city, "updated_at": time.Now()},
	}

	var user User
	err := s.usersCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetAllUsers(ctx context.Context) ([]*User, error) {
	cursor, err := s.usersCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := cursor.Close(ctx); err != nil {
			slog.Error("Couldn't close the mongoDB connection correctly", "error", err)
		}
	}()

	var users []*User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) Close(ctx context.Context) error {
	disconnectCtx, disconnectCancel := context.WithTimeout(ctx, 5*time.Second)
	defer disconnectCancel()

	return s.client.Disconnect(disconnectCtx)
}
