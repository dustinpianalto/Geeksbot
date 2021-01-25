package database

import (
	"database/sql"
	"log"

	"github.com/dustinpianalto/geeksbot"
)

type channelService struct {
	db *sql.DB
}

func (s channelService) Channel(id string) (geeksbot.Channel, error) {
	var channel geeksbot.Channel
	var guild_id string
	queryString := "SELECT id, guild_id, admin, default_channel, new_patron FROM channels WHERE id = $1"
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&channel.ID, &guild_id, &channel.Admin, &channel.Default, &channel.NewPatron)
	if err != nil {
		return geeksbot.Channel{}, err
	}
	guild, err := GuildService.Guild(guild_id)
	if err != nil {
		return geeksbot.Channel{}, err
	}
	channel.Guild = guild
	return channel, nil
}

func (s channelService) CreateChannel(c geeksbot.Channel) (geeksbot.Channel, error) {
	queryString := "INSERT INTO channels (id, guild_id, admin, default_channel, new_patron) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(queryString, c.ID, c.Guild.ID, c.Admin, c.Default, c.NewPatron)
	return c, err
}

func (s channelService) DeleteChannel(c geeksbot.Channel) error {
	queryString := "DELETE FROM channels WHERE id = $1"
	_, err := s.db.Exec(queryString, c.ID)
	return err
}

func (s channelService) UpdateChannel(c geeksbot.Channel) (geeksbot.Channel, error) {
	queryString := "UPDATE channels SET admin = $2, default_channel = $3, new_patron = $4 WHERE id = $1"
	_, err := s.db.Exec(queryString, c.ID, c.Admin, c.Default, c.NewPatron)
	return c, err
}

func (s channelService) GuildChannels(g geeksbot.Guild) ([]geeksbot.Channel, error) {
	var channels []geeksbot.Channel
	queryString := "SELECT id FROM channels WHERE guild_id = $1"
	rows, err := s.db.Query(queryString, g.ID)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		channel, err := s.Channel(id)
		if err != nil {
			log.Println(err)
			continue
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (s channelService) GetOrCreateChannel(id string, guild_id string) (geeksbot.Channel, error) {
	channel, err := s.Channel(id)
	if err.Error() == sql.ErrNoRows.Error() {
		var guild geeksbot.Guild
		guild, err = GuildService.GetOrCreateGuild(guild_id)
		if err != nil {
			return geeksbot.Channel{}, err
		}
		channel, err = s.CreateChannel(geeksbot.Channel{
			ID:        id,
			Guild:     guild,
			Admin:     false,
			Default:   false,
			NewPatron: false,
		})
	}
	return channel, err
}
