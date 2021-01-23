package geeksbot

type User struct {
	ID       string
	SteamID  string
	IsActive bool
	IsStaff  bool
	IsAdmin  bool
}

type UserService interface {
	User(id string) (User, error)
	CreateUser(u User) (User, error)
	DeleteUser(u User) error
	UpdateUser(u User) (User, error)
}
