package message

import (
	"context"
)

type MessageSender interface {
	SendMessage(text string, UserID int64) error
}

type Model struct {
	ctx      context.Context
	tgClient MessageSender // Client.
}

type Message struct {
	Text            string
	UserID          int64
	UserName        string
	UserDisplayName string
}

func New(ctx context.Context, tgClient MessageSender) *Model {
	return &Model{
		ctx:      ctx,
		tgClient: tgClient,
	}
}

func (s *Model) IncomingMessage(msg Message) error {
	//TODO: some logic
	return nil
}
