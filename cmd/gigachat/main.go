package main

import (
	"github.com/KillReall666/Antispam-tg-bot/internal/config/gigachatcfg"
	"github.com/KillReall666/Antispam-tg-bot/internal/service/gigachat"
	"github.com/KillReall666/Antispam-tg-bot/internal/storage/redis"
	"log"
	"net/http"
)

func main() {
	cfg, err := gigachatcfg.New()
	if err != nil {
		log.Fatal("error loading gigachatcfg config", err)
	}

	redisClient := redis.NewRedisClient(cfg.RedisAddr)
	pong, err := redisClient.Ping()
	if err != nil {
		log.Fatal("redis connection error:", err)
	}
	log.Println("connection to redis established:", pong)

	service := gigachat.New(cfg, redisClient)
	//service.GetAccessTokenAuthRequests()
	//service.GetRequest("userName")
	//service.Consumer()
	service.HuggingFaceAPIRequest()

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting http server", err)
	}

}
