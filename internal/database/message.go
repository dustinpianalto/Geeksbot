package database

import (
	"database/sql"
	"log"

	"github.com/dustinpianalto/geeksbot"
)

type messageService struct {
	db *sql.DB
}

func (s messageService) Message(id string) (geeksbot.Message, error) {
	var m geeksbot.Message
	var channel_id string
	var author_id string
	queryString := `SELECT m.id, m.created_at, m.modified_at, m.content, 
						m.previous_content, m.channel_id, m.author_id
						FROM messages as m
						WHERE m.id = $1`
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&m.ID, &m.CreatedAt, &m.ModifiedAt, &m.Content,
		&m.PreviousContent, &channel_id, &author_id)
	if err != nil {
		return geeksbot.Message{}, err
	}
	author, err := UserService.User(author_id)
	if err != nil {
		return geeksbot.Message{}, err
	}
	m.Author = author
	channel, err := ChannelService.Channel(channel_id)
	if err != nil {
		return geeksbot.Message{}, err
	}
	m.Channel = channel
	return m, nil
}

func (s messageService) CreateMessage(m geeksbot.Message) (geeksbot.Message, error) {
	queryString := `INSERT INTO messages (id, created_at, content, channel_id, author_id) 
					VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(queryString, m.ID, m.CreatedAt, m.Content, m.Channel.ID, m.Author.ID)
	return m, err
}

func (s messageService) UpdateMessage(m geeksbot.Message) (geeksbot.Message, error) {
	var content string
	var previousContent []string
	queryString := "SELECT content, previous_content FROM messages WHERE id = $1"
	row := s.db.QueryRow(queryString, m.ID)
	err := row.Scan(&content, &previousContent)
	if err != nil {
		return geeksbot.Message{}, err
	}
	if m.Content != content {
		previousContent = append(previousContent, content)
	}
	queryString = "UPDATE messages SET modified_at = $2, content = $3, previous_content = $4 WHERE id = $1"
	_, err = s.db.Exec(queryString, m.ID, m.ModifiedAt, m.Content, previousContent)
	m.PreviousContent = previousContent
	return m, nil
}

func (s messageService) DeleteMessage(m geeksbot.Message) error {
	queryString := "DELETE FROM messages WHERE id = $1"
	_, err := s.db.Exec(queryString, m.ID)
	return err
}

func (s messageService) ChannelMessages(c geeksbot.Channel) ([]geeksbot.Message, error) {
	var messages []geeksbot.Message
	queryString := `SELECT id FROM messages WHERE channel_id = $1`
	rows, err := s.db.Query(queryString, c.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		m, err := s.Message(id)
		if err != nil {
			log.Println(err)
			continue
		}
		messages = append(messages, m)
	}
	return messages, nil
}
