package gigachat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/KillReall666/Antispam-tg-bot/internal/config/gigachatcfg"
	"github.com/KillReall666/Antispam-tg-bot/internal/model/gigachat"
	"github.com/KillReall666/Antispam-tg-bot/internal/storage/redis"
	"github.com/hupe1980/go-huggingface"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"os"
	"time"

	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
)

// Тут нам надо подключить не саму БД а интерфейс типа UserRepository repository.UserRepository  (в New тоже самое) а сам интерфейс описать на уровне storage
type Service struct {
	cfg       *gigachatcfg.GigachatConfig
	client    *http.Client
	redis     *redis.RedisClient
	token     *string
	expiresAr *int64
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func New(cfg *gigachatcfg.GigachatConfig, redis *redis.RedisClient) *Service {
	client := http.Client{}
	return &Service{
		cfg:    cfg,
		client: &client,
		redis:  redis,
	}
}

// Consumer Receiving messages from tg bot
func (s *Service) Consumer() {
	env, err := stream.NewEnvironment(stream.NewEnvironmentOptions())
	if err != nil {
		log.Fatalf("Failed to create environment: %v", err)
	}

	streamName := "app-gigachat-stream"
	env.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		},
	)

	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		log.Printf("Stream: %s - Received message: %s\n", consumerContext.Consumer.GetStreamName(),
			message.Data)
	}

	consumer, err := env.NewConsumer(streamName, messagesHandler,
		stream.NewConsumerOptions().SetOffset(stream.OffsetSpecification{}.First()))
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	log.Println(" [x] Waiting for messages. enter to close the consumer")
	_, _ = reader.ReadString('\n')
	err = consumer.Close()

}

func (s *Service) GetAccessTokenAuthRequests() {
	ctx := context.Background()
	payload := strings.NewReader("scope=" + s.cfg.Scope)
	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.AuthUrl+OAuthPath, payload)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("RqUID", uuid.NewString())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s.cfg.ClientID+":"+s.cfg.ClientSecret)))

	resp, err := s.client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("unexpected status code %d %v", resp.StatusCode, resp.Status)
	}

	var oauth OAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&oauth)
	if err != nil {
		log.Println(err)
	}

	s.token = &oauth.AccessToken
	s.expiresAr = &oauth.ExpiresAt

	expVal := *s.expiresAr
	expTime := time.Duration(expVal) * time.Nanosecond

	err = s.redis.Set("userName", *s.token, expTime)
	if err != nil {
		log.Println("err to set redis key: ", err)
	}
}

func (s *Service) Print() {
	fmt.Println(":)")
}

func (s *Service) GetRequest(userName string) {
	tkn, err := s.redis.Get(userName)
	if err != nil {
		log.Printf("err getting token for user %s: %v", userName, err)
	}

	if tkn == "" {
		tkn = *s.token
	}

	chatModel := gigachat.ChatModel{
		Model: "GigaChat",
		Messages: []gigachat.Message{
			{
				Role:    "user",
				Content: "Проверка на мат! Содержит ли следующиая строка матные слова: Застрахуй команду корабля со скипидаром? Просто ответь да или нет",
			},
		},
		Temperature:       1.0,
		TopP:              0.1,
		N:                 1,
		Stream:            false,
		MaxTokens:         512,
		RepetitionPenalty: 1,
	}

	reqBody, err := json.Marshal(chatModel)
	if err != nil {
		log.Println("err when marshal json: ", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", s.cfg.RequestUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("err when making req", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tkn)

	res, err := client.Do(req)
	if err != nil {
		log.Println("err when doing req", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("err when reading res body", err)
		return
	}

	var completion gigachat.Completion

	err = json.Unmarshal(body, &completion)
	if err != nil {
		log.Println("err when unmarshal:", err)
		return
	}

	//fmt.Println(completion.Choices.Message)
}

func (s *Service) HuggingFaceAPIRequest() {
	ic := huggingface.NewInferenceClient("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")

	res, err := ic.TextGeneration(context.Background(), &huggingface.TextGenerationRequest{
		Inputs: "Can you speak russian?",
		Model:  "BlinkDL/rwkv-5-world",
	})

	if err != nil {
		log.Println("hugging api err: ", err)
	}

	fmt.Println(res[0].GeneratedText)
}
