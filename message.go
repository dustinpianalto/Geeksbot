package geeksbot

import "github.com/bwmarrin/discordgo"

type Message struct {
	ID              string
	CreatedAt       int64
	ModifiedAt      int64
	Content         string
	PreviousContent []string
	Channel         *Channel
	Author          *User
	Embed           *discordgo.MessageEmbed
	PreviousEmbeds  []*discordgo.MessageEmbed
}

type MessageService interface {
	Message(id string) (*Message, error)
	CreateMessage(m *Message) (*Message, error)
	DeleteMessage(m *Message) error
	ChannelMessages(c *Channel) ([]*Message, error)
	UpdateMessage(m *Message) (*Message, error)
}
