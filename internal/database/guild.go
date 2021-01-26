package database

import (
	"database/sql"
	"log"

	"github.com/dustinpianalto/geeksbot"
	"github.com/lib/pq"
)

type guildService struct {
	db *sql.DB
}

func (s guildService) Guild(id string) (geeksbot.Guild, error) {
	var g geeksbot.Guild
	queryString := "SELECT id, new_patron_message, prefixes FROM guilds WHERE id = $1"
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&g.ID, &g.NewPatronMessage, pq.Array(&g.Prefixes))
	if err != nil {
		return geeksbot.Guild{}, err
	}
	return g, nil
}

func (s guildService) CreateGuild(g geeksbot.Guild) (geeksbot.Guild, error) {
	queryString := "INSERT INTO guilds (id, new_patron_message, prefixes) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(queryString, g.ID, g.NewPatronMessage, pq.Array(g.Prefixes))
	return g, err
}

func (s guildService) DeleteGuild(g geeksbot.Guild) error {
	queryString := "DELETE FROM guilds WHERE id = $1"
	_, err := s.db.Exec(queryString, g.ID)
	return err
}

func (s guildService) UpdateGuild(g geeksbot.Guild) (geeksbot.Guild, error) {
	queryString := "UPDATE guilds SET new_patron_message = $2, prefixes = $3 WHERE id = $1"
	_, err := s.db.Exec(queryString, g.ID, g.NewPatronMessage, pq.Array(g.Prefixes))
	return g, err
}

func (s guildService) GuildRoles(g geeksbot.Guild) ([]geeksbot.Role, error) {
	var roles []geeksbot.Role
	queryString := "SELECT id FROM roles WHERE guild_id = $1"
	rows, err := s.db.Query(queryString, g.ID)
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
		role, err := s.Role(id)
		if err != nil {
			log.Println(err)
			continue
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (s guildService) CreateRole(r geeksbot.Role) (geeksbot.Role, error) {
	queryString := "INSERT INTO roles (id, role_type, guild_id) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(queryString, r.ID, r.RoleType, r.Guild.ID)
	return r, err
}

func (s guildService) Role(id string) (geeksbot.Role, error) {
	var role geeksbot.Role
	var guild_id string
	queryString := "SELECT id, role_type, guild_id FROM roles WHERE id = $1"
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&role.ID, &role.RoleType, &guild_id)
	if err != nil {
		return geeksbot.Role{}, err
	}
	guild, err := s.Guild(guild_id)
	if err != nil {
		return geeksbot.Role{}, err
	}
	role.Guild = guild
	return role, nil
}

func (s guildService) UpdateRole(r geeksbot.Role) (geeksbot.Role, error) {
	queryString := "UPDATE roles SET role_type = $2 WHERE id = $1"
	_, err := s.db.Exec(queryString, r.ID, r.RoleType)
	return r, err
}

func (s guildService) DeleteRole(r geeksbot.Role) error {
	queryString := "DELETE FROM roles WHERE id = $1"
	_, err := s.db.Exec(queryString, r.ID)
	return err
}

func (s guildService) GetOrCreateGuild(id string) (geeksbot.Guild, error) {
	guild, err := s.Guild(id)
	if err == sql.ErrNoRows {
		guild = geeksbot.Guild{
			ID:       id,
			Prefixes: []string{},
		}
		guild, err = s.CreateGuild(guild)
	}
	return guild, err
}

func (s guildService) CreateOrUpdateRole(r geeksbot.Role) (geeksbot.Role, error) {
	role, err := s.Role(r.ID)
	if err.Error() == `pq: duplicate key value violates unique constraint "roles_pkey"` {
		role, err = s.UpdateRole(r)
	}
	return role, err
}
