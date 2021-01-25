package geeksbot

import "database/sql"

type User struct {
	ID       string
	SteamID  sql.NullString
	IsActive bool
	IsStaff  bool
	IsAdmin  bool
}

type UserService interface {
	User(id string) (User, error)
	CreateUser(u User) (User, error)
	DeleteUser(u User) error
	UpdateUser(u User) (User, error)
	GetOrCreateUser(id string) (User, error)
}
