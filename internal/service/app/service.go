package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KillReall666/Antispam-tg-bot/internal/config/appcfg"
	"github.com/KillReall666/Antispam-tg-bot/internal/model/app"
	"log"
	"net/http"
	"strings"
)

type MessageSender interface {
	SendMessageToBot(text string, UserID int64) error
	SendMessageToGroup(text string, groupID int64) error
	//DeleteMessageFromGroup(groupID int64, messageID int) error
}

type UserUnderAttackStorage interface {
	//DB METHODS
}

type Service struct {
	ctx      context.Context
	tgClient MessageSender // Клиент.
	storage  UserUnderAttackStorage
	cfg      *appcfg.AppConfig
}

func New(ctx context.Context, tgClient MessageSender, userUnderAttack UserUnderAttackStorage, cfg *appcfg.AppConfig) *Service {
	return &Service{
		ctx:      ctx,
		tgClient: tgClient,
		storage:  userUnderAttack,
		cfg:      cfg,
	}
}

func (s *Service) IncomingMessage(msg app.Message) error {
	if strings.Contains(msg.Text, "колбаса") {
		err := s.DeleteMessageFromGroup(msg.GroupID, msg.MessageID)
		log.Println("err deleted message from group:", err)
	}
	return nil

}

func (s *Service) DeleteMessageFromGroup(groupID int64, messageID int) error {
	requestBody := &app.DeleteMessageRequestBody{
		ChatID:    groupID,
		MessageID: messageID,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling DeleteMessageRequestBody: %v", err))
	}

	request, err := http.NewRequest("POST", s.cfg.TgApiURL+s.cfg.Token+DeleteMsg, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return errors.New(fmt.Sprintf("error creating request: %v", err))
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("error doing request: %v", err))
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed to delete message: %v", response.Status))
	}

	return nil
}

//func (s *Service) SendQuestionToGigaChatSber(text string, UserID int64) error {
//	req, err := http.NewRequest("POST", "https://gigachat.devices.sberbank.ru/api/v1/models", []byte(text)
//0}
