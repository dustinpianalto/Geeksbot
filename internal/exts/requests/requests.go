package requests

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/discord_utils"
	"github.com/dustinpianalto/geeksbot/pkg/services"
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
						ctx.Message.Author.Mention(),
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

var CloseCommand = &disgoman.Command{
	Name:                "close",
	Aliases:             nil,
	Description:         "Close a request and mark it as completed.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              closeCommandFunc,
}

func closeCommandFunc(ctx disgoman.Context, args []string) {
	var ids []int64
	var reason string
	_, err := strconv.ParseInt(args[len(args)-1], 10, 64)
	if err != nil {
		reason = args[len(args)-1]
		args = args[0 : len(args)-1]
	}
	for _, a := range args {
		a = strings.Trim(a, ",")
		id, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%s is not a valid request id", a), err)
			continue
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		discord_utils.SendErrorMessage(ctx, "No requests to close", nil)
		return
	}
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error getting Guild from the database", err)
		return
	}
	closer, err := services.UserService.GetOrCreateUser(ctx.Message.Author.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error closing the request. Could not get user.", err)
		return
	}
	int64ID, _ := strconv.ParseInt(ctx.Message.ID, 10, 64)
	s := discord_utils.ParseSnowflake(int64ID)
	for _, id := range ids {
		request, err := services.RequestService.Request(id)
		if err != nil || request.Guild.ID != guild.ID {
			discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%d is not a request in this guild.", id), err)
			continue
		}
		if request.Author != closer && !closer.IsStaff && !closer.IsAdmin {
			if !discord_utils.IsGuildMod(ctx, closer) || !discord_utils.IsGuildAdmin(ctx, closer) {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("You are not authorized to close %d", id), nil)
				continue
			}
		}
		if request.Completed {
			discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%d is already closed", id), err)
			continue
		}
		request.Completed = true
		request.CompletedAt.Valid = true
		request.CompletedAt.Time = s.CreationTime
		request.CompletedBy = &closer
		if reason != "" {
			request.CompletedMessage.Valid = true
			request.CompletedMessage.String = reason
		}
		request, err = services.RequestService.UpdateRequest(request)
		if err != nil {
			discord_utils.SendErrorMessage(ctx, fmt.Sprintf("Error closing %d", id), err)
			continue
		}
		_, err = ctx.Send(fmt.Sprintf("%d has been closed", request.ID))
		if err != nil {
			discord_utils.SendErrorMessage(ctx, "There was an error sending the message. The request was closed.", err)
		}

	}
}

var ListCommand = &disgoman.Command{
	Name:                "list",
	Aliases:             []string{"request_list", "requests_list"},
	Description:         "List your open requests or all open requests if the caller is a moderator",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              listCommandFunc,
}

func listCommandFunc(ctx disgoman.Context, args []string) {
	user, err := services.UserService.GetOrCreateUser(ctx.Message.Author.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error closing the request. Could not get user.", err)
		return
	}
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error getting Guild from the database", err)
		return
	}

	var requests []geeksbot.Request
	if discord_utils.IsGuildMod(ctx, user) || discord_utils.IsGuildAdmin(ctx, user) {
		requests, err = services.RequestService.GuildRequests(guild, false)
	} else {
		requests, err = services.RequestService.UserRequests(user, false)
	}

	for _, request := range requests {
		var authorName string
		author, err := ctx.Session.GuildMember(guild.ID, request.Author.ID)
		if err != nil {
			authorName = "Unknown"
		} else {
			if author.Nick == "" {
				authorName = author.User.Username
			} else {
				authorName = author.Nick
			}
		}

		var channelName string
		channel, err := ctx.Session.Channel(request.Channel.ID)
		if err != nil {
			channelName = "Unknown"
		} else {
			channelName = channel.Name
		}

		_, _ = ctx.Send(fmt.Sprintf("```md\n"+
			"< Request ID            Requested By >\n"+
			"< %11s %-23s >\n"+
			"%s\n\n"+
			"Comments: Not Implemented Yet\n"+
			"Requested At: %s\n"+
			"In: %s\n"+
			"```",
			request.ID,
			authorName,
			request.Content,
			request.RequestedAt.Format("2006-01-02 15:04:05 MST"),
			channelName,
		))
	}
}
