package services

import (
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/pkg/database"
)

var (
	GuildService   geeksbot.GuildService
	UserService    geeksbot.UserService
	ChannelService geeksbot.ChannelService
	MessageService geeksbot.MessageService
	PatreonService geeksbot.PatreonService
	RequestService geeksbot.RequestService
	ServerService  geeksbot.ServerService
)

func InitializeServices() {
	GuildService = database.GuildService
	UserService = database.UserService
	ChannelService = database.ChannelService
	MessageService = database.MessageService
	PatreonService = database.PatreonService
	RequestService = database.RequestService
	ServerService = database.ServerService
}
