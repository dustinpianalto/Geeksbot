package geeksbot

import (
	"database/sql"
	"time"
)

type Request struct {
	ID               int64
	Author           User
	Channel          Channel
	Guild            Guild
	Content          string
	RequestedAt      time.Time
	Completed        bool
	CompletedAt      sql.NullTime
	CompletedBy      *User
	Message          Message
	CompletedMessage sql.NullString
}

type Comment struct {
	ID        int64
	Author    User
	Request   Request
	CommentAt time.Time
	Content   string
}

type RequestService interface {
	Request(id int64) (Request, error)
	UserRequests(u User, completed bool) ([]Request, error)
	GuildRequests(g Guild, completed bool) ([]Request, error)
	CreateRequest(r Request) (Request, error)
	UpdateRequest(r Request) (Request, error)
	DeleteRequest(r Request) error
	Comment(id int64) (Comment, error)
	RequestComments(r Request) ([]Comment, error)
	RequestCommentCount(r Request) (int, error)
	CreateComment(c Comment) (Comment, error)
	DeleteComment(c Comment) error
}
