package requests

import (
	"strconv"
	"strings"

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
		discord_utils.SendErrorMessage(ctx, "Error creating the request", err)
		return
	}
	channel, err := services.ChannelService.GetOrCreateChannel(ctx.Message.ChannelID, ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error creating the request", err)
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
}
