package geeksbot

type Server struct {
	ID              int
	Name            string
	IPAddr          string
	Port            int
	Password        string
	AlertsChannel   *Channel
	Guild           Guild
	InfoChannel     *Channel
	InfoMessage     *Message
	SettingsMessage *Message
}

type ServerService interface {
	ServerByID(id int) (Server, error)
	ServerByName(name string, guild Guild) (Server, error)
	CreateServer(s Server) (Server, error)
	DeleteServer(s Server) error
	UpdateServer(s Server) (Server, error)
	GuildServers(g Guild) ([]Server, error)
}
