package guild

import (
	"fmt"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/discord_utils"
	"github.com/dustinpianalto/geeksbot/internal/services"
)

var AddPrefixCommand = &disgoman.Command{
	Name:                "addPrefix",
	Aliases:             []string{"ap"},
	Description:         "Add a prefix for use in this guild.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              addPrefixCommandFunc,
}

func addPrefixCommandFunc(ctx disgoman.Context, args []string) {
	if len(args) == 0 {
		discord_utils.SendErrorMessage(ctx, "Please include at least one prefix to add", nil)
		return
	}
	guild, err := services.GuildService.Guild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Guild not configured, adding new guild to the database", nil)
		guild = geeksbot.Guild{
			ID:       ctx.Guild.ID,
			Prefixes: args,
		}
		guild, err = services.GuildService.CreateGuild(guild)
		if err != nil {
			discord_utils.SendErrorMessage(ctx, "Error adding guild to database.", err)
			return
		}
	} else {
		guild.Prefixes = append(guild.Prefixes, args...)
		guild, err = services.GuildService.UpdateGuild(guild)
		if err != nil {
			discord_utils.SendErrorMessage(ctx, "Error adding prefixes to guild.", err)
			return
		}
	}
	_, err = ctx.Send(fmt.Sprintf("Prefixes Updates.\nThe Prefixes for this guild are currently %#v", guild.Prefixes))
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error sending update message", err)
	}
}

var RemovePrefixCommand = &disgoman.Command{
	Name:                "removePrefix",
	Aliases:             []string{"rp"},
	Description:         "Remove a prefix so it can't be used in this guild",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              removePrefixCommandFunc,
}

func removePrefixCommandFunc(ctx disgoman.Context, args []string) {
	if len(args) == 0 {
		discord_utils.SendErrorMessage(ctx, "Please include at least one prefix to remove", nil)
		return
	}
	guild, err := services.GuildService.Guild(ctx.Guild.ID)
	var removed []string
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Guild not configured, adding new guild to the database", nil)
		guild = geeksbot.Guild{
			ID:       ctx.Guild.ID,
			Prefixes: []string{},
		}
		guild, err = services.GuildService.CreateGuild(guild)
		if err != nil {
			discord_utils.SendErrorMessage(ctx, "Error adding guild to database.", err)
			return
		}
	} else {
		for _, a := range args {
			l := len(guild.Prefixes) - 1
			for i := 0; i < l; i++ {
				if a == guild.Prefixes[i] {
					guild.Prefixes = append(guild.Prefixes[:i], guild.Prefixes[i+1:]...)
					l--
					i--
					removed = append(removed, a)
				}
			}
		}
		guild, err = services.GuildService.UpdateGuild(guild)
		if err != nil {
			discord_utils.SendErrorMessage(ctx, "Error removing prefixes from guild.", err)
			return
		}
	}
	_, err = ctx.Send(fmt.Sprintf("Prefixes Updates.\n"+
		"The Prefixes for this guild are currently %#v\n"+
		"Removed: %#v", guild.Prefixes, removed))
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Error sending update message", err)
	}
}
