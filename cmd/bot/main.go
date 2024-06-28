package main

import (
	"context"
	config "github.com/KillReall666/Antispam-tg-bot/internal/config/appcfg"
	service "github.com/KillReall666/Antispam-tg-bot/internal/service/app"
	"github.com/KillReall666/Antispam-tg-bot/internal/storage/postgres"
	"log"

	tg "github.com/KillReall666/Antispam-tg-bot/internal/client"
)

func main() {

	log.Println("Starting Telegram Bot")

	ctx := context.Background()

	appCfg, err := config.New()
	if err != nil {
		log.Fatal("error loading appcfg config", err)
	}

	tgProcessingFuncHandler := tg.HandlerFunc(tg.ProcessingMessages)

	//Инициализация DB.
	db, err := postgres.New(appCfg.DBConnStr)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	//Инициализация tg client.
	tgClient, err := tg.New(appCfg, tgProcessingFuncHandler)
	if err != nil {
		log.Fatal("error initialization tgClient", err)
	}

	service := service.New(ctx, tgClient, db, appCfg)

	tgClient.ListenUpdates(service)

}
