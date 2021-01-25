package requests

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/discord_utils"
	"github.com/dustinpianalto/geeksbot/internal/services"
)

var RequestCommand = &disgoman.Command{
	Name:                "request",
	Aliases:             nil,
	Description:         "Submit a request for the guild staff",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              requestCommandFunc,
}

func requestCommandFunc(ctx disgoman.Context, args []string) {
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error getting Guild from the database", err)
		return
	}
	requestMsg := strings.Join(args, " ")
	author, err := services.UserService.GetOrCreateUser(ctx.Message.Author.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error creating the request. Could not get user.", err)
		return
	}
	channel, err := services.ChannelService.GetOrCreateChannel(ctx.Message.ChannelID, ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error creating the request. Could not get channel.", err)
		return
	}
	int64ID, _ := strconv.ParseInt(ctx.Message.ID, 10, 64)
	s := discord_utils.ParseSnowflake(int64ID)
	message, err := services.MessageService.CreateMessage(geeksbot.Message{
		ID:        ctx.Message.ID,
		CreatedAt: s.CreationTime,
		Content:   ctx.Message.Content,
		Channel:   channel,
		Author:    author,
	})
	request := geeksbot.Request{
		Author:      author,
		Channel:     channel,
		Guild:       guild,
		Content:     requestMsg,
		RequestedAt: s.CreationTime,
		Completed:   false,
		Message:     message,
	}
	request, err = services.RequestService.CreateRequest(request)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error creating the request", err)
		return
	}
	channels, err := services.ChannelService.GuildChannels(guild)
	if err == nil {
		var mentionRolesString string
		roles, err := services.GuildService.GuildRoles(guild)
		if err == nil {
			for _, r := range roles {
				if r.RoleType == "admin" || r.RoleType == "moderator" {
					mentionRolesString += fmt.Sprintf("<@&%s> ", r.ID)
				}
			}
		}
		for _, c := range channels {
			if c.Admin {
				_, _ = ctx.Session.ChannelMessageSend(c.ID,
					fmt.Sprintf("%s\n"+
						"New Request ID %d "+
						"%s has requested assistance: \n"+
						"```\n%s\n```\n"+
						"Requested At: %s\n"+
						"In: %s",
						mentionRolesString,
						request.ID,
						ctx.Message.Author.Mention,
						request.Content,
						request.RequestedAt.UTC().Format(time.UnixDate),
						ctx.Channel.Mention(),
					),
				)
			}
		}
	}
	_, err = ctx.Send(fmt.Sprintf("%s The admin have recieved your request.\n "+
		"If you would like to close or add a comment to this request please reference ID `%v`",
		ctx.Message.Author.Mention(), request.ID,
	))
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "There was an error sending the message. The request was created.", err)
	}
}
