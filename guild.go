package geeksbot

import "database/sql"

type Guild struct {
	ID               string
	NewPatronMessage sql.NullString
	Prefixes         []string
}

type Role struct {
	ID       string
	RoleType string
	Guild    Guild
}

type GuildService interface {
	Guild(id string) (Guild, error)
	CreateGuild(g Guild) (Guild, error)
	DeleteGuild(g Guild) error
	UpdateGuild(g Guild) (Guild, error)
	GuildRoles(g Guild) ([]Role, error)
	CreateRole(r Role) (Role, error)
	Role(id string) (Role, error)
	UpdateRole(r Role) (Role, error)
	DeleteRole(r Role) error
}
