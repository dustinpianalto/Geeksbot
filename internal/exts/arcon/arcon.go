package arcon

import (
	"fmt"
	"strings"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/discord_utils"
	"github.com/dustinpianalto/geeksbot/pkg/services"
	"github.com/gorcon/rcon"
)

var ListplayersCommand = &disgoman.Command{
	Name:                "request",
	Aliases:             nil,
	Description:         "Submit a request for the guild staff",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              listplayersCommandFunc,
}

func listplayersCommandFunc(ctx disgoman.Context, args []string) {
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error getting Guild from the database", err)
		return
	}
	author, err := services.UserService.GetOrCreateUser(ctx.Message.Author.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Sorry, there was a problem getting your user.", err)
		return
	}
	if !discord_utils.IsGuildAdmin(ctx, author) && !discord_utils.IsGuildMod(ctx, author) {
		return
	}
	if len(args) == 0 {
		servers, err := services.ServerService.GuildServers(guild)
		if err != nil {
			discord_utils.SendErrorMessage(ctx, "Could not find any servers for this guild", err)
			return
		}
		for _, server := range servers {
			go listplayers(ctx, server)
		}
		return
	}
	serverName := strings.Join(args, " ")
	server, err := services.ServerService.ServerByName(serverName, guild)
	if err != nil {
		discord_utils.SendErrorMessage(ctx,
			fmt.Sprintf("Could not find **%s** in this guild.", serverName),
			err,
		)
		return
	}
	listplayers(ctx, server)
}

func listplayers(ctx disgoman.Context, server geeksbot.Server) {
	msg, err := ctx.Send(fmt.Sprintf("**Getting data for %s**", server.Name))
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "There was an error getting the player list", err)
		return
	}
	conn, err := rcon.Dial(fmt.Sprintf("%s:%d", server.IPAddr, server.Port), server.Password)
	if err != nil {
		_, _ = ctx.Session.ChannelMessageEdit(ctx.Channel.ID, msg.ID,
			fmt.Sprintf("**Could not open connection to %s**", server.Name),
		)
		return
	}
	defer conn.Close()
	response, err := conn.Execute("listplayers")
	if err != nil {
		_, _ = ctx.Session.ChannelMessageEdit(ctx.Channel.ID, msg.ID,
			fmt.Sprintf("**There was a problem getting a response from %s**", server.Name),
		)
		return
	}
	if strings.HasPrefix(response, "No Players") {
		_, _ = ctx.Session.ChannelMessageEdit(ctx.Channel.ID, msg.ID,
			fmt.Sprintf("**%s: %s**", server.Name, response),
		)
		return
	}
	players := strings.Split(response, "\n")
	for i, player := range players {
		parts := strings.Split(player, ", ")
		steamID := parts[len(parts)-1]
		user, err := services.UserService.GetBySteamID(steamID)
		if err == nil {
			duser, err := ctx.Session.GuildMember(ctx.Guild.ID, user.ID)
			if err == nil {
				players[i] = fmt.Sprintf("%s (%s)", player, duser.Mention())
			}
		}
	}
	_, _ = ctx.Session.ChannelMessageEdit(ctx.Channel.ID, msg.ID,
		fmt.Sprintf("**%s:**\n%s", server.Name, strings.Join(players, "\n")),
	)
}
