package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
)

func main() {
	Token := os.Getenv("DISCORDGO_TOKEN")
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

	//postgres.ConnectDatabase(os.Getenv("DATABASE_URL"))
	//postgres.InitializeDatabase()
	//utils.LoadTestData()

	//us := &postgres.UserService{DB: postgres.DB}
	//gs := &postgres.GuildService{DB: postgres.DB}

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

	// Add Command Handlers
	exts.AddCommandHandlers(&manager)
	//services.InitalizeServices(us, gs)

	//if _, ok := handler.Commands["help"]; !ok {
	//	handler.AddDefaultHelpCommand()
	//}

	dg.AddHandler(manager.OnMessage)
	dg.AddHandler(manager.StatusManager.OnReady)
	//dg.AddHandler(guild_management.OnMessageUpdate)
	//dg.AddHandler(guild_management.OnMessageDelete)
	//dg.AddHandler(user_management.OnGuildMemberAddLogging)
	//dg.AddHandler(user_management.OnGuildMemberRemoveLogging)

	err = dg.Open()
	if err != nil {
		log.Println("There was an error opening the connection, ", err)
		return
	}

	// Start the Error handler in a goroutine
	go ErrorHandler(manager.ErrorChannel)

	// Start the Logging handler in a goroutine
	//go logging.LoggingHandler(logging.LoggingChannel)

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
