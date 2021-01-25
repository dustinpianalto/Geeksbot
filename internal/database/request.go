package database

import (
	"database/sql"
	"log"

	"github.com/dustinpianalto/geeksbot"
)

type requestService struct {
	db *sql.DB
}

func (s requestService) Request(id int64) (geeksbot.Request, error) {
	var r geeksbot.Request
	var aID string
	var cID string
	var gID string
	var uID sql.NullString
	var mID string
	queryString := `SELECT id, author_id, channel_id, guild_id, content, requested_at, completed,
						completed_at, completed_by, message_id, completed_message
						WHERE id = $1`
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&r.ID, &aID, &cID, &gID, &r.Content, &r.RequestedAt, &r.Completed,
		&r.CompletedAt, &uID, &mID, &r.CompletedMessage)
	if err != nil {
		return geeksbot.Request{}, err
	}
	author, err := UserService.User(aID)
	if err != nil {
		return geeksbot.Request{}, err
	}
	guild, err := GuildService.Guild(gID)
	if err != nil {
		return geeksbot.Request{}, err
	}
	channel, err := ChannelService.Channel(cID)
	if err != nil {
		return geeksbot.Request{}, err
	}
	if !uID.Valid {
		r.CompletedBy = nil
	} else {
		completedBy, err := UserService.User(uID.String)
		if err != nil {
			return geeksbot.Request{}, err
		}
		r.CompletedBy = &completedBy
	}
	message, err := MessageService.Message(mID)
	if err != nil {
		return geeksbot.Request{}, err
	}
	r.Author = author
	r.Guild = guild
	r.Channel = channel
	r.Message = message
	return r, nil
}

func (s requestService) UserRequests(u geeksbot.User, completed bool) ([]geeksbot.Request, error) {
	var requests []geeksbot.Request
	var queryString string
	if completed {
		queryString = "SELECT id FROM requests WHERE author_id = $1"
	} else {
		queryString = "SELECT id FROM requests WHERE author_id = $1 AND completed = False"
	}
	rows, err := s.db.Query(queryString, u.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		request, err := s.Request(id)
		if err != nil {
			log.Println(err)
			continue
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func (s requestService) GuildRequests(g geeksbot.Guild, completed bool) ([]geeksbot.Request, error) {
	var requests []geeksbot.Request
	var queryString string
	if completed {
		queryString = "SELECT id FROM requests WHERE guild_id = $1"
	} else {
		queryString = "SELECT id FROM requests WHERE guild_id = $1 AND completed = False"
	}
	rows, err := s.db.Query(queryString, g.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		request, err := s.Request(id)
		if err != nil {
			log.Println(err)
			continue
		}
		requests = append(requests, request)
	}
	return requests, nil

}

func (s requestService) CreateRequest(r geeksbot.Request) (geeksbot.Request, error) {
	queryString := `INSERT INTO requests 
						(author_id, channel_id, guild_id, content, requested_at, 
							completed, completed_at, completed_by, message_id, completed_message)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	var id int64
	var completedByID sql.NullString
	if r.CompletedBy == nil {
		completedByID.String = ""
		completedByID.Valid = false
	} else {
		completedByID.String = r.CompletedBy.ID
		completedByID.Valid = true
	}
	err := s.db.QueryRow(queryString,
		r.Author.ID,
		r.Channel.ID,
		r.Guild.ID,
		r.Content,
		r.RequestedAt,
		r.Completed,
		r.CompletedAt,
		completedByID,
		r.Message.ID,
		r.CompletedMessage).Scan(&id)
	if err != nil {
		return geeksbot.Request{}, err
	}
	r.ID = id
	return r, nil
}

func (s requestService) UpdateRequest(r geeksbot.Request) (geeksbot.Request, error) {
	queryString := `UPDATE requests SET 
						completed = $2, completed_at = $3, completed_by = $4, completed_message = $5
						WHERE id = $1`
	_, err := s.db.Exec(queryString, r.ID, r.Completed, r.CompletedAt, r.CompletedBy, r.CompletedMessage)
	return r, err
}

func (s requestService) DeleteRequest(r geeksbot.Request) error {
	queryString := "DELETE FROM requests WHERE id = $1"
	_, err := s.db.Exec(queryString, r.ID)
	return err
}

func (s requestService) Comment(id int64) (geeksbot.Comment, error) {
	var c geeksbot.Comment
	var aID string
	var rID int64
	queryString := "SELECT id, author_id, request_id, comment_at, content WHERE id = $1"
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&c.ID, &aID, &rID, &c.CommentAt, &c.Content)
	if err != nil {
		return geeksbot.Comment{}, err
	}
	author, err := UserService.User(aID)
	if err != nil {
		return geeksbot.Comment{}, err
	}
	c.Author = author
	request, err := s.Request(rID)
	if err != nil {
		return geeksbot.Comment{}, err
	}
	c.Request = request
	return c, nil
}

func (s requestService) RequestComments(r geeksbot.Request) ([]geeksbot.Comment, error) {
	var comments []geeksbot.Comment
	queryString := "SELECT id FROM comments WHERE request_id = $1"
	rows, err := s.db.Query(queryString, r.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		comment, err := s.Comment(id)
		if err != nil {
			log.Println(err)
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (s requestService) RequestCommentCount(r geeksbot.Request) (int, error) {
	var count int
	queryString := "SELECT COUNT(id) FROM comments WHERE request_id = $1"
	row := s.db.QueryRow(queryString, r.ID)
	err := row.Scan(&count)
	return count, err
}

func (s requestService) CreateComment(c geeksbot.Comment) (geeksbot.Comment, error) {
	queryString := `INSERT INTO comments (author_id, request_id, comment_at, content)
						VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := s.db.QueryRow(queryString, c.Author.ID, c.Request.ID, c.CommentAt, c.Content).Scan(&id)
	if err != nil {
		return geeksbot.Comment{}, err
	}
	c.ID = id
	return c, nil
}

func (s requestService) DeleteComment(c geeksbot.Comment) error {
	queryString := "DELETE FROM comments WHERE id = $1"
	_, err := s.db.Exec(queryString, c.ID)
	return err
}
