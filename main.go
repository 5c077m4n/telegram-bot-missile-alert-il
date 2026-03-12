package main

import (
	"context"
	"os"
	"os/signal"

	"telegram-bot-missile-alert-il/telebot"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := telebot.Setup(ctx); err != nil {
		panic(err)
	}

}
