package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/database"
	"github.com/dustinpianalto/geeksbot/internal/exts"
)

func main() {
	Token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("There was an error when creating the Discord Session, ", err)
		return
	}
	dg.State.MaxMessageCount = 100
	dg.StateEnabled = true

	dg.Identify = discordgo.Identify{
		Intents: discordgo.MakeIntent(discordgo.IntentsAll),
	}

	database.ConnectDatabase(os.Getenv("DATABASE_URL"))

	owners := []string{
		"351794468870946827",
	}

	manager := disgoman.CommandManager{
		Prefixes:         getPrefixes,
		Owners:           owners,
		StatusManager:    disgoman.GetDefaultStatusManager(),
		ErrorChannel:     make(chan disgoman.CommandError, 10),
		Commands:         make(map[string]*disgoman.Command),
		IgnoreBots:       true,
		CheckPermissions: false,
	}

	geeksbot := geeksbot.Geeksbot{
		GuildService:   database.GuildService,
		UserService:    database.UserService,
		ChannelService: database.ChannelService,
		MessageService: database.MessageService,
		PatreonService: database.PatreonService,
		RequestService: database.RequestService,
		ServerService:  database.ServerService,
		CommandManager: manager,
	}

	// Add Command Handlers
	exts.AddCommandHandlers(&geeksbot)

	dg.AddHandler(geeksbot.OnMessage)
	dg.AddHandler(geeksbot.StatusManager.OnReady)

	err = dg.Open()
	if err != nil {
		log.Println("There was an error opening the connection, ", err)
		return
	}

	// Start the Error handler in a goroutine
	go ErrorHandler(geeksbot.ErrorChannel)

	log.Println("The Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Shutting Down...")
	err = dg.Close()
	if err != nil {
		log.Println(err)
	}
}

func getPrefixes(guildID string) []string {
	return []string{"G.", "g."}
}

func ErrorHandler(ErrorChan chan disgoman.CommandError) {
	for ce := range ErrorChan {
		msg := ce.Message
		if msg == "" {
			msg = ce.Error.Error()
		}
		_, _ = ce.Context.Send(msg)
		log.Println(ce.Error)
	}
}
