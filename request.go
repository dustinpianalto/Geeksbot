package geeksbot

import "time"

type Request struct {
	ID               int64
	Author           *User
	Channel          *Channel
	Content          string
	RequestedAt      time.Time
	CompletedAt      time.Time
	CompletedBy      *User
	Message          *Message
	CompletedMessage string
}

type Comment struct {
	ID        int64
	Author    *User
	Request   *Request
	CommentAt time.Time
	Content   string
}

type RequestService interface {
	Request(id int64) (*Request, error)
	UserRequests(u *User, completed bool) ([]*Request, error)
	GuildRequests(g *Guild, completed bool) ([]*Request, error)
	CreateRequest(r *Request) (*Request, error)
	UpdateRequest(r *Request) (*Request, error)
	DeleteRequest(r *Request) error
	Comment(id int64) (*Comment, error)
	RequestComments(r *Request) ([]*Comment, error)
	RequestCommentCount(r *Request) (int, error)
	CreateComment(c *Comment) (*Comment, error)
	DeleteComment(c *Comment) error
}