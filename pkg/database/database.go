package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dustinpianalto/geeksbot/pkg/database/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

var (
	db             *sql.DB
	GuildService   guildService
	UserService    userService
	ChannelService channelService
	MessageService messageService
	RequestService requestService
	PatreonService patreonService
	ServerService  serverService
)

func ConnectDatabase(dbConnString string) {
	var err error
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatal(fmt.Errorf("Can't connect to the database: %w", err))
	}
	log.Println("Database Connected.")
	db.SetMaxOpenConns(75)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(300)
	initServices()
	log.Println("Services Initialized")
}

func RunMigrations() {
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	d, err := bindata.WithInstance(s)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot load migrations: %w", err))
	}
	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(fmt.Errorf("cannot create db driver: %w", err))
	}
	m, err := migrate.NewWithInstance("go-bindatafoo", d, "postgres", instance)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot create migration instance: %w", err))
	}
	err = m.Up()
	if err != nil {
		if err.Error() == "no change" {
			log.Println(err)
		} else {
			log.Fatal(err)
		}
	}
	log.Println("Migrations Run")

}

func initServices() {
	GuildService = guildService{db: db}
	UserService = userService{db: db}
	ChannelService = channelService{db: db}
	MessageService = messageService{db: db}
	PatreonService = patreonService{db: db}
	RequestService = requestService{db: db}
	ServerService = serverService{db: db}
}
