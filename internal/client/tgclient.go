package tg

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/KillReall666/Antispam-tg-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(tgUpdate tgbotapi.Update, c *Client, msgModel *message.Model)

func (f HandlerFunc) RunFunc(tgUpdate tgbotapi.Update, c *Client, msgModel *message.Model) {
	f(tgUpdate, c, msgModel)
}

type Client struct {
	client               *tgbotapi.BotAPI
	handleProcessingFunc HandlerFunc //Функция обработки входящих сообщений.
}

type TokenGetter interface {
	Token() string
}

func New(tokenGetter TokenGetter, handleProcessingFunc HandlerFunc) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("NewBotAPI error: %v", err))
	}
	return &Client{
		client:               client,
		handleProcessingFunc: handleProcessingFunc,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ParseMode = "markdown"
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.New(fmt.Sprintf("Error msg send client.Send: %v", err))
	}
	return nil
}

func (c *Client) ListenUpdates(msgModel *message.Model) {
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
func ProcessingMessages(tgUpdate tgbotapi.Update, c *Client, msgModel *message.Model) {
	if tgUpdate.Message != nil {
		//Пользователь написал сообщение
		log.Println(fmt.Sprintf("[%s][%v] %s", tgUpdate.Message.From.UserName, tgUpdate.Message.From.ID, tgUpdate.Message.Text))
		err := msgModel.IncomingMessage(message.Message{
			Text:            tgUpdate.Message.Text,
			UserID:          tgUpdate.Message.From.ID,
			UserName:        tgUpdate.Message.From.UserName,
			UserDisplayName: strings.TrimSpace(tgUpdate.Message.From.FirstName + " " + tgUpdate.Message.From.LastName),
		})
		if err != nil {
			log.Println(fmt.Sprintf("error processing message: %v", err))
		}
	}
}
