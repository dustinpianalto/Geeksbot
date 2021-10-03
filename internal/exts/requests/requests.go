package requests

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/discord_utils"
	"github.com/dustinpianalto/geeksbot/internal/utils"
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
						"New Request ID %d\n"+
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
			if !discord_utils.IsGuildMod(ctx, closer) && !discord_utils.IsGuildAdmin(ctx, closer) {
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
		dmChannel, err := ctx.Session.UserChannelCreate(request.Author.ID)
		if err != nil {
			return
		}
		_, _ = ctx.Session.ChannelMessageSend(dmChannel.ID,
			fmt.Sprintf("%s has closed request %d which you opened in the %s channel.\n```%s```\n",
				discord_utils.GetDisplayName(ctx, request.CompletedBy.ID),
				request.ID,
				discord_utils.GetChannelName(ctx, request.Channel.ID),
				request.Content,
			))
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
		commentCount, err := services.RequestService.RequestCommentCount(request)
		if err != nil {
			commentCount = 0
		}
		_, _ = ctx.Send(fmt.Sprintf("```md\n"+
			"< Request ID             Requested By >\n"+
			"< %-11d %23s >\n"+
			"%s\n\n"+
			"Comments: %d\n"+
			"Requested At: %s\n"+
			"In: %s\n"+
			"```",
			request.ID,
			authorName,
			request.Content,
			commentCount,
			request.RequestedAt.Format("2006-01-02 15:04:05 MST"),
			channelName,
		))
	}

	_, _ = ctx.Send(fmt.Sprintf("```There %s currently %d open %s```", utils.PluralizeString("is", len(requests)), len(requests), utils.PluralizeString("request", len(requests))))
}

var CommentCommand = &disgoman.Command{
	Name:                "comment",
	Aliases:             []string{"update", "add_comment"},
	Description:         "Add a comment to an existing request.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              commentCommandFunc,
}

func commentCommandFunc(ctx disgoman.Context, args []string) {
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error getting Guild from the database", err)
		return
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Please include the ID of the request to update.", err)
		return
	}
	message := strings.Join(args[1:len(args)], " ")
	author, err := services.UserService.GetOrCreateUser(ctx.Message.Author.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Sorry, there was an issue finding your user account", err)
		return
	}
	request, err := services.RequestService.Request(id)
	if err != nil || request.Guild.ID != guild.ID {
		discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%d is not a valid request for this guild", id), err)
		return
	}
	int64ID, _ := strconv.ParseInt(ctx.Message.ID, 10, 64)
	s := discord_utils.ParseSnowflake(int64ID)
	comment := geeksbot.Comment{
		Author:    author,
		Request:   request,
		CommentAt: s.CreationTime,
		Content:   message,
	}
	comment, err = services.RequestService.CreateComment(comment)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "There was a problem adding your comment", err)
		return
	}
	channels, err := services.ChannelService.GuildChannels(guild)
	if err == nil {
		comments, _ := services.RequestService.RequestComments(request)
		var commentString string
		var commentStrings []string
		commentString = fmt.Sprintf("Comment added:\n```md\n"+
			"< Request ID             Requested By >\n"+
			"< %-11d %23s >\n"+
			"%s\n\n"+
			"Comments: Not Implemented Yet\n"+
			"Requested At: %s\n"+
			"In: %s\n"+
			"```",
			request.ID,
			discord_utils.GetDisplayName(ctx, request.Author.ID),
			request.Content,
			request.RequestedAt.Format("2006-01-02 15:04:05"),
			discord_utils.GetChannelName(ctx, request.Channel.ID),
		)
		for _, c := range comments {
			if err != nil {
				log.Println(err)
				continue
			}
			cs := fmt.Sprintf("```md\n%s\n- %s At %s\n```\n",
				c.Content,
				discord_utils.GetDisplayName(ctx, c.Author.ID),
				c.CommentAt.Format("2006-01-02 15:04:05"),
			)
			if len(commentString+cs) >= 2000 {
				commentStrings = append(commentStrings, commentString)
				commentString = ""
			}
			commentString += cs
		}
		commentStrings = append(commentStrings, commentString)
		for _, c := range channels {
			if c.Admin {
				for _, s := range commentStrings {
					_, _ = ctx.Session.ChannelMessageSend(c.ID, s)
				}
			}
		}
	} else {
		log.Println(err)
	}
	_, err = ctx.Send(fmt.Sprintf("%s your comment has been added.", ctx.Message.Author.Mention()))
	dmChannel, err := ctx.Session.UserChannelCreate(request.Author.ID)
	if err != nil {
		return
	}
	_, _ = ctx.Session.ChannelMessageSend(dmChannel.ID,
		fmt.Sprintf("%s has add a comment to request %d which you opened in the %s channel.\n```%s```\n```%s```",
			discord_utils.GetDisplayName(ctx, author.ID),
			request.ID,
			discord_utils.GetChannelName(ctx, ctx.Channel.ID),
			request.Content,
			message,
		))

}

var ViewCommand = &disgoman.Command{
	Name:                "view",
	Aliases:             nil,
	Description:         "View the details about a request.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              viewCommandFunc,
}

func viewCommandFunc(ctx disgoman.Context, args []string) {
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error getting Guild from the database", err)
		return
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Please include the ID of the request to view.", err)
		return
	}
	request, err := services.RequestService.Request(id)
	if err != nil || request.Guild.ID != guild.ID {
		discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%d is not a valid request in this guild", id), err)
		return
	}
	requestor, err := services.UserService.GetOrCreateUser(ctx.Message.Author.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Sorry, there was an issue finding your user account", err)
		return
	}
	if request.Author.ID != ctx.Message.Author.ID &&
		!discord_utils.IsGuildMod(ctx, requestor) &&
		!discord_utils.IsGuildAdmin(ctx, requestor) {
		discord_utils.SendErrorMessage(ctx, "You are not authorized to view that request", nil)
		return
	}
	comments, err := services.RequestService.RequestComments(request)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "There was an error getting the comments.", err)
	}
	var commentString string
	var commentStrings []string
	commentString = fmt.Sprintf("```md\n"+
		"< Request ID             Requested By >\n"+
		"< %-11d %23s >\n"+
		"%s\n\n"+
		"Requested At: %s\n"+
		"In: %s\n"+
		"```",
		request.ID,
		discord_utils.GetDisplayName(ctx, request.Author.ID),
		request.Content,
		request.RequestedAt.Format("2006-01-02 15:04:05"),
		discord_utils.GetChannelName(ctx, request.Channel.ID),
	)
	for _, c := range comments {
		cs := fmt.Sprintf("```md\n%s\n- %s At %s\n```\n",
			c.Content,
			discord_utils.GetDisplayName(ctx, c.Author.ID),
			c.CommentAt.Format("2006-01-02 15:04:05"),
		)
		if len(commentString+cs) >= 2000 {
			commentStrings = append(commentStrings, commentString)
			commentString = ""
		}
		commentString += cs
	}
	commentStrings = append(commentStrings, commentString)
	for _, c := range commentStrings {
		_, _ = ctx.Send(c)
	}
}
