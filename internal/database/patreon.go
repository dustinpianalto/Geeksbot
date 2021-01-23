package database

import (
	"database/sql"
	"log"

	"github.com/dustinpianalto/geeksbot"
)

type patreonService struct {
	db *sql.DB
}

func (s patreonService) PatreonCreatorByID(id int) (geeksbot.PatreonCreator, error) {
	var creator geeksbot.PatreonCreator
	var gID string
	queryString := "SELECT id, creator, link, guild_id FROM patreon_creator WHERE id = $1"
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&creator.ID, &creator.Creator, &creator.Link, &gID)
	if err != nil {
		return geeksbot.PatreonCreator{}, err
	}
	guild, err := GuildService.Guild(gID)
	if err != nil {
		return geeksbot.PatreonCreator{}, err
	}
	creator.Guild = guild
	return creator, nil
}

func (s patreonService) PatreonCreatorByName(name string, g geeksbot.Guild) (geeksbot.PatreonCreator, error) {
	var id int
	queryString := "SELECT id FROM patreon_creator WHERE creator = $1 AND guild_id = $2"
	err := s.db.QueryRow(queryString, name, g.ID).Scan(&id)
	if err != nil {
		return geeksbot.PatreonCreator{}, nil
	}
	creator, err := s.PatreonCreatorByID(id)
	return creator, err
}

func (s patreonService) CreatePatreonCreator(c geeksbot.PatreonCreator) (geeksbot.PatreonCreator, error) {
	var id int
	queryString := `INSERT INTO patreon_creator (creator, link, guild_id) VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRow(queryString, c.Creator, c.Link, c.Guild.ID).Scan(&id)
	if err != nil {
		return geeksbot.PatreonCreator{}, err
	}
	c.ID = id
	return c, nil
}

func (s patreonService) UpdatePatreonCreator(c geeksbot.PatreonCreator) (geeksbot.PatreonCreator, error) {
	queryString := `UPDATE patreon_creator SET creator = $2, link = $3 WHERE id = $1`
	_, err := s.db.Exec(queryString, c.ID, c.Creator, c.Link)
	return c, err
}

func (s patreonService) DeletePatreonCreator(c geeksbot.PatreonCreator) error {
	queryString := `DELETE FROM patreon_creator WHERE id = $1`
	_, err := s.db.Exec(queryString, c.ID)
	return err
}

func (s patreonService) PatreonTierByID(id int) (geeksbot.PatreonTier, error) {
	var tier geeksbot.PatreonTier
	var cID int
	var rID string
	var next int
	queryString := `SELECT id, name, description, creator, role, next_tier FROM patreon_tier WHERE id = id`
	err := s.db.QueryRow(queryString, id).Scan(&tier.ID, &tier.Name, &tier.Description, &cID, &rID, &next)
	if err != nil {
		return geeksbot.PatreonTier{}, err
	}
	creator, err := s.PatreonCreatorByID(cID)
	if err != nil {
		return geeksbot.PatreonTier{}, err
	}
	tier.Creator = creator
	role, err := GuildService.Role(rID)
	if err != nil {
		return geeksbot.PatreonTier{}, err
	}
	tier.Role = role
	if next == -1 {
		tier.NextTier = nil
		return tier, nil
	}
	nextTier, err := s.PatreonTierByID(next)
	if err != nil {
		return geeksbot.PatreonTier{}, err
	}
	tier.NextTier = &nextTier
	return tier, nil
}

func (s patreonService) PatreonTierByName(name string, creator string) (geeksbot.PatreonTier, error) {
	var id int
	queryString := `SELECT id FROM patreon_tier WHERE name = $1 AND creator = $2`
	err := s.db.QueryRow(queryString, name, creator).Scan(&id)
	if err != nil {
		return geeksbot.PatreonTier{}, err
	}
	tier, err := s.PatreonTierByID(id)
	return tier, err
}

func (s patreonService) CreatePatreonTier(t geeksbot.PatreonTier) (geeksbot.PatreonTier, error) {
	var id int
	queryString := `INSERT INTO patreon_tier (name, description, creator, role, next_tier)
						VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.db.QueryRow(queryString, t.Name, t.Description, t.Creator.ID, t.Role.ID, t.NextTier.ID).Scan(&id)
	if err != nil {
		return geeksbot.PatreonTier{}, err
	}
	t.ID = id
	return t, nil
}

func (s patreonService) UpdatePatreonTier(t geeksbot.PatreonTier) (geeksbot.PatreonTier, error) {
	queryString := `UPDATE patreon_tier SET name = $2, description = $3, role = $4, next_tier = $5 WHERE id = $1`
	_, err := s.db.Exec(queryString, t.ID, t.Name, t.Description, t.Role.ID, t.NextTier.ID)
	return t, err
}

func (s patreonService) DeletePatreonTier(t geeksbot.PatreonTier) error {
	queryString := `DELETE FROM patreon_tier WHERE id = $1`
	_, err := s.db.Exec(queryString, t.ID)
	return err
}

func (s patreonService) GuildPatreonCreators(g geeksbot.Guild) ([]geeksbot.PatreonCreator, error) {
	var creators []geeksbot.PatreonCreator
	queryString := `SELECT id FROM patreon_creator WHERE guild_id = $1`
	rows, err := s.db.Query(queryString, g.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		creator, err := s.PatreonCreatorByID(id)
		if err != nil {
			log.Println(err)
			continue
		}
		creators = append(creators, creator)
	}
	return creators, nil
}

func (s patreonService) CreatorPatreonTiers(c geeksbot.PatreonCreator) ([]geeksbot.PatreonTier, error) {
	var tiers []geeksbot.PatreonTier
	queryString := `SELECT id FROM patreon_tier WHERE creator = $1`
	rows, err := s.db.Query(queryString, c.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		tier, err := s.PatreonTierByID(id)
		if err != nil {
			log.Println(err)
			continue
		}
		tiers = append(tiers, tier)
	}
	return tiers, nil
}
