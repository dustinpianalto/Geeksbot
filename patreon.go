package geeksbot

type PatreonCreator struct {
	ID      int
	Creator string
	Link    string
	Guild   Guild
}

type PatreonTier struct {
	ID          int
	Name        string
	Description string
	Creator     PatreonCreator
	Role        Role
	NextTier    *PatreonTier
}

type PatreonService interface {
	PatreonCreatorByID(id int) (PatreonCreator, error)
	PatreonCreatorByName(name string, guild Guild) (PatreonCreator, error)
	CreatePatreonCreator(c PatreonCreator) (PatreonCreator, error)
	UpdatePatreonCreator(c PatreonCreator) (PatreonCreator, error)
	DeletePatreonCreator(c PatreonCreator) error
	PatreonTierByID(id int) (PatreonTier, error)
	PatreonTierByName(name string, creator string) (PatreonTier, error)
	CreatePatreonTier(t PatreonTier) (PatreonTier, error)
	UpdatePatreonTier(t PatreonTier) (PatreonTier, error)
	DeletePatreonTier(t PatreonTier) error
	GuildPatreonCreators(g Guild) ([]PatreonCreator, error)
	CreatorPatreonTiers(c PatreonCreator) ([]PatreonTier, error)
}
