package app

type Message struct {
	Text            string
	UserID          int64
	UserName        string
	UserDisplayName string
	GroupID         int64
	MessageID       int
}

type DeleteMessageRequestBody struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
}
