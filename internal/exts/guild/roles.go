package guild

import (
	"fmt"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/discord_utils"
	"github.com/dustinpianalto/geeksbot/internal/services"
)

var AddModeratorRoleCommand = &disgoman.Command{
	Name:                "addMod",
	Aliases:             []string{"addModerator", "addModRole"},
	Description:         "Add a role which is allowed to run moderator commands",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              addModeratorRoleCommandFunc,
}

func addModeratorRoleCommandFunc(ctx disgoman.Context, args []string) {
	var count int
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Something went wrong getting the guild", err)
		return
	}
	roles := append(args, ctx.Message.MentionRoles...)
	if len(roles) > 0 {
		for _, id := range ctx.Message.MentionRoles {
			_, err := services.GuildService.CreateRole(geeksbot.Role{
				ID:       id,
				RoleType: "moderator",
				Guild:    guild,
			})
			if err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("There was a problem adding <@&%s>", id), err)
				continue
			}
			count++
			_, _ = ctx.Send(fmt.Sprintf("Added <@&%s> as a moderator role.", id))
		}
	} else {
		_, _ = ctx.Send("Please include at least one role to make a moderator role.")
	}
}
