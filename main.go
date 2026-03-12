package main

import (
	"context"
	"os"
	"os/signal"

	"telegram-bot-missile-alert-il/telebot"
	"telegram-bot-missile-alert-il/users"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	store, err := users.NewStore(uri, dbName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := store.Close(ctx); err != nil {
			panic(err)
		}
	}()

	if err := telebot.Setup(ctx, store); err != nil {
		panic(err)
	}
}
