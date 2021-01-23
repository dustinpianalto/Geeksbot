package geeksbot

import "github.com/dustinpianalto/disgoman"

type Geeksbot struct {
	GuildService   GuildService
	UserService    UserService
	ChannelService ChannelService
	MessageService MessageService
	RequestService RequestService
	PatreonService PatreonService
	ServerService  ServerService
	disgoman.CommandManager
}
