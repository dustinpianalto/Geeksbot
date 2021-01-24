package database

import (
	"database/sql"
	"log"

	"github.com/dustinpianalto/geeksbot"
)

type serverService struct {
	db *sql.DB
}

func (s serverService) ServerByID(id int) (geeksbot.Server, error) {
	var server geeksbot.Server
	var guildID string
	var aChanID sql.NullString
	var iChanID sql.NullString
	var iMsgID sql.NullString
	var sMsgID sql.NullString
	queryString := `SELECT id, name, ip_address, port, password, alerts_channel_id, 
						guild_id, info_channel_id, info_message_id, settings_message_id
						FROM servers WHERE id = $1`
	row := s.db.QueryRow(queryString, id)
	err := row.Scan(&server.ID, &server.Name, &server.IPAddr, &server.Port, &server.Password,
		&aChanID, &guildID, &iChanID, &iMsgID, &sMsgID)
	if err != nil {
		return geeksbot.Server{}, err
	}
	guild, err := GuildService.Guild(guildID)
	if err != nil {
		return geeksbot.Server{}, err
	}
	if !aChanID.Valid {
		server.AlertsChannel = nil
	} else {
		alertChannel, err := ChannelService.Channel(aChanID.String)
		if err != nil {
			return geeksbot.Server{}, err
		}
		server.AlertsChannel = &alertChannel
	}
	if !iChanID.Valid {
		server.InfoChannel = nil
	} else {
		infoChannel, err := ChannelService.Channel(iChanID.String)
		if err != nil {
			return geeksbot.Server{}, err
		}
		server.InfoChannel = &infoChannel
	}
	if !iMsgID.Valid {
		server.InfoMessage = nil
	} else {
		infoMessage, err := MessageService.Message(iMsgID.String)
		if err != nil {
			return geeksbot.Server{}, err
		}
		server.InfoMessage = &infoMessage
	}
	if !sMsgID.Valid {
		server.SettingsMessage = nil
	} else {
		settingsMessage, err := MessageService.Message(sMsgID.String)
		if err != nil {
			return geeksbot.Server{}, err
		}
		server.SettingsMessage = &settingsMessage
	}
	server.Guild = guild
	return server, nil
}

func (s serverService) ServerByName(name string, guild geeksbot.Guild) (geeksbot.Server, error) {
	var id int
	queryString := "SELECT id FROM servers WHERE name = $1 AND guild_id = $2"
	row := s.db.QueryRow(queryString, name, guild.ID)
	err := row.Scan(&id)
	if err != nil {
		return geeksbot.Server{}, err
	}
	server, err := s.ServerByID(id)
	return server, err
}

func (s serverService) CreateServer(server geeksbot.Server) (geeksbot.Server, error) {
	var id int
	queryString := `INSERT INTO servers (name, ip_address, port, password, alerts_channel_id, 
						guild_id, info_channel_id, info_message_id, settings_message_id)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	err := s.db.QueryRow(queryString,
		server.Name,
		server.IPAddr,
		server.Port,
		server.Password,
		server.AlertsChannel,
		server.Guild,
		server.InfoChannel,
		server.InfoMessage,
		server.SettingsMessage,
	).Scan(&id)
	if err != nil {
		return geeksbot.Server{}, err
	}
	server.ID = id
	return server, nil
}

func (s serverService) DeleteServer(server geeksbot.Server) error {
	queryString := `DELETE FROM servers WHERE id = $1`
	_, err := s.db.Exec(queryString, server.ID)
	return err
}

func (s serverService) UpdateServer(server geeksbot.Server) (geeksbot.Server, error) {
	queryString := `UPDATE servers SET name = $2, ip_address = $3, port = $4, password = $5,
						alerts_channel_id = $6, info_channel_id = $7, info_message_id = $8,
						settings_message_id = $9 WHERE id = $1`
	_, err := s.db.Exec(queryString,
		server.Name,
		server.IPAddr,
		server.Port,
		server.Password,
		server.AlertsChannel.ID,
		server.InfoChannel.ID,
		server.InfoMessage.ID,
		server.SettingsMessage.ID,
	)
	return server, err
}

func (s serverService) GuildServers(g geeksbot.Guild) ([]geeksbot.Server, error) {
	var servers []geeksbot.Server
	queryString := `SELECT id FROM servers WHERE guild_id = $1`
	rows, err := s.db.Query(queryString, g.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		server, err := s.ServerByID(id)
		if err != nil {
			log.Println(err)
			continue
		}
		servers = append(servers, server)
	}
	return servers, nil
}
