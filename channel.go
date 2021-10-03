package geeksbot

type Channel struct {
	ID        string
	Guild     Guild
	Admin     bool
	Default   bool
	NewPatron bool
}

type ChannelService interface {
	Channel(id string) (Channel, error)
	CreateChannel(c Channel) (Channel, error)
	DeleteChannel(c Channel) error
	GuildChannels(g Guild) ([]Channel, error)
	UpdateChannel(c Channel) (Channel, error)
	GetOrCreateChannel(id string, guild_id string) (Channel, error)
}
