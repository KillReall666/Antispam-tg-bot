package main

import (
	"context"
	"log"

	tg "github.com/KillReall666/Antispam-tg-bot/internal/client"
	config "github.com/KillReall666/Antispam-tg-bot/internal/config"
	message "github.com/KillReall666/Antispam-tg-bot/internal/model"
)

func main() {

	log.Println("Starting Telegram Bot")

	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("error loading config", err)
	}

	tgProcessingFuncHandler := tg.HandlerFunc(tg.ProcessingMessages)

	//Инициализация tg client.
	tgClient, err := tg.New(cfg, tgProcessingFuncHandler)
	if err != nil {
		log.Fatal("error initialization tgClient", err)
	}

	msgModel := message.New(ctx, tgClient)
	tgClient.ListenUpdates(msgModel)

}
