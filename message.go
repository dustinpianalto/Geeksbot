package geeksbot

import (
	"database/sql"
	"time"
)

type Message struct {
	ID              string
	CreatedAt       time.Time
	ModifiedAt      sql.NullTime
	Content         string
	PreviousContent []string
	Channel         Channel
	Author          User
}

type MessageService interface {
	Message(id string) (Message, error)
	CreateMessage(m Message) (Message, error)
	DeleteMessage(m Message) error
	ChannelMessages(c Channel) ([]Message, error)
	UpdateMessage(m Message) (Message, error)
}
