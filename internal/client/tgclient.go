package tg

import (
	"errors"
	"fmt"
	app2 "github.com/KillReall666/Antispam-tg-bot/internal/model/app"
	"github.com/KillReall666/Antispam-tg-bot/internal/service/app"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Добавить коммит 4
type HandlerFunc func(tgUpdate tgbotapi.Update, c *Client, msgModel *app.Service)

func (f HandlerFunc) RunFunc(tgUpdate tgbotapi.Update, c *Client, msgModel *app.Service) {
	f(tgUpdate, c, msgModel)
}

type Client struct {
	client               *tgbotapi.BotAPI
	handleProcessingFunc HandlerFunc //Функция обработки входящих сообщений.
}

type TokenGetter interface {
	GetToken() string
}

func New(tokenGetter TokenGetter, handleProcessingFunc HandlerFunc) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.GetToken())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("newBotAPI error: %v", err))
	}
	return &Client{
		client:               client,
		handleProcessingFunc: handleProcessingFunc,
	}, nil
}

func (c *Client) SendMessageToBot(text string, userID int64) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ParseMode = "markdown"
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.New(fmt.Sprintf("error msg send client.Send: %v", err))
	}
	return nil
}

func (c *Client) SendMessageToGroup(text string, groupID int64) error {
	msg := tgbotapi.NewMessage(groupID, text)
	msg.ParseMode = "markdown"
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.New(fmt.Sprintf("error msg send client.Send: %v", err))
	}
	return nil
}

func (c *Client) ListenUpdates(msgModel *app.Service) {
	updateConf := tgbotapi.NewUpdate(0)
	updateConf.Timeout = 60

	updates := c.client.GetUpdatesChan(updateConf)

	log.Println("Start listening for tg messages...")
	for update := range updates {
		//Функция обработки сообщений обернутая в middleware.
		c.handleProcessingFunc.RunFunc(update, c, msgModel)
	}
}

// ProcessingMessages функция обработки сообщений
func ProcessingMessages(tgUpdate tgbotapi.Update, c *Client, msgModel *app.Service) {
	if tgUpdate.Message != nil {
		//Пользователь написал сообщение
		log.Println(fmt.Sprintf("userName:[%s] userID:[%v] chatID:[%v] messageID:[%v]: %s", tgUpdate.Message.From.UserName, tgUpdate.Message.From.ID, tgUpdate.Message.Chat.ID, tgUpdate.Message.MessageID, tgUpdate.Message.Text))
		err := msgModel.Producer(app2.Message{
			Text:            tgUpdate.Message.Text,
			UserID:          tgUpdate.Message.From.ID,
			UserName:        tgUpdate.Message.From.UserName,
			UserDisplayName: strings.TrimSpace(tgUpdate.Message.From.FirstName + " " + tgUpdate.Message.From.LastName),
			GroupID:         tgUpdate.Message.Chat.ID,
			MessageID:       tgUpdate.Message.MessageID,
		})
		if err != nil {
			log.Println(fmt.Sprintf("error processing message: %v", err))
		}
	}
}
