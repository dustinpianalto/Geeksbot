package database

import (
	"database/sql"

	"github.com/dustinpianalto/geeksbot"
)

type userService struct {
	db *sql.DB
}

func (s userService) User(id string) (geeksbot.User, error) {
	var user geeksbot.User
	queryString := "SELECT id, steam_id, active, staff, admin WHERE id = $1"
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&user.ID, &user.SteamID, &user.IsActive, &user.IsStaff, &user.IsAdmin)
	return user, err
}

func (s userService) CreateUser(u geeksbot.User) (geeksbot.User, error) {
	queryString := "INSERT INTO users (id, steam_id, active, staff, admin) VALUES ($1, $2, $3, $4, $5)"
	var err error
	if u.SteamID.Valid {
		_, err = s.db.Exec(queryString, u.ID, u.SteamID.String, u.IsActive, u.IsStaff, u.IsAdmin)
	} else {
		_, err = s.db.Exec(queryString, u.ID, nil, u.IsActive, u.IsStaff, u.IsAdmin)
	}
	return u, err
}

func (s userService) DeleteUser(u geeksbot.User) error {
	queryString := "DELETE FROM users WHERE id = $1"
	_, err := s.db.Exec(queryString, u.ID)
	return err
}

func (s userService) UpdateUser(u geeksbot.User) (geeksbot.User, error) {
	queryString := "UPDATE users SET steam_id = $2, active = $3, staff = $4, admin = $5 WHERE id = $1"
	_, err := s.db.Exec(queryString, u.ID, u.SteamID, u.IsActive, u.IsStaff, u.IsAdmin)
	return u, err
}
