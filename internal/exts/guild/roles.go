package guild

import (
	"fmt"
	"strings"

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
	added := make(map[string]bool)
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Something went wrong getting the guild", err)
		return
	}
	roles := append(args, ctx.Message.MentionRoles...)
	if len(roles) > 0 {
		for _, id := range roles {
			if strings.HasPrefix(id, "<@&") && strings.HasSuffix(id, ">") {
				continue
			}
			if _, ok := added[id]; ok {
				continue
			}
			if _, err = ctx.Session.State.Role(ctx.Guild.ID, id); err != nil {
				_, _ = ctx.Send(fmt.Sprintf("%s does not reference a valid role for this guild.", id))
				continue
			}
			_, err := services.GuildService.CreateOrUpdateRole(geeksbot.Role{
				ID:       id,
				RoleType: "moderator",
				Guild:    guild,
			})
			if err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("There was a problem adding <@&%s>", id), err)
				continue
			}
			added[id] = true
			count++
			_, _ = ctx.Send(fmt.Sprintf("Added <@&%s> as a moderator role.", id))
		}
		_, _ = ctx.Send(fmt.Sprintf("Added %d moderator roles.", count))
	} else {
		_, _ = ctx.Send("Please include at least one role to make a moderator role.")
	}
}

var AddAdminRoleCommand = &disgoman.Command{
	Name:                "addAdmin",
	Aliases:             []string{"addAdminRole"},
	Description:         "Add a role which is allowed to run admin commands",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              addAdminRoleCommandFunc,
}

func addAdminRoleCommandFunc(ctx disgoman.Context, args []string) {
	var count int
	added := make(map[string]bool)
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Something went wrong getting the guild", err)
		return
	}
	roles := append(args, ctx.Message.MentionRoles...)
	if len(roles) > 0 {
		for _, id := range roles {
			if strings.HasPrefix(id, "<@&") && strings.HasSuffix(id, ">") {
				continue
			}
			if _, ok := added[id]; ok {
				continue
			}
			if _, err = ctx.Session.State.Role(ctx.Guild.ID, id); err != nil {
				_, _ = ctx.Send(fmt.Sprintf("%s does not reference a valid role for this guild.", id))
				continue
			}
			_, err := services.GuildService.CreateOrUpdateRole(geeksbot.Role{
				ID:       id,
				RoleType: "admin",
				Guild:    guild,
			})
			if err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("There was a problem adding <@&%s>", id), err)
				continue
			}
			added[id] = true
			count++
			_, _ = ctx.Send(fmt.Sprintf("Added <@&%s> as an admin role.", id))
		}
		_, _ = ctx.Send(fmt.Sprintf("Added %d admin roles.", count))
	} else {
		_, _ = ctx.Send("Please include at least one role to make an admin role.")
	}
}

var RemoveModRoleCommand = &disgoman.Command{
	Name:                "removeMod",
	Aliases:             []string{"removeModeratorRole", "removeModRole", "removeAdmin", "removeAdminRole"},
	Description:         "Remove a role or several roles from the moderator or admin list",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              removeModRoleCommandFunc,
}

func removeModRoleCommandFunc(ctx disgoman.Context, args []string) {
	var count int
	added := make(map[string]bool)
	guild, err := services.GuildService.GetOrCreateGuild(ctx.Guild.ID)
	if err != nil {
		discord_utils.SendErrorMessage(ctx, "Something went wrong getting the guild", err)
		return
	}
	roles := append(args, ctx.Message.MentionRoles...)
	if len(roles) > 0 {
		for _, id := range roles {
			if strings.HasPrefix(id, "<@&") && strings.HasSuffix(id, ">") {
				continue
			}
			if _, ok := added[id]; ok {
				continue
			}
			if _, err = ctx.Session.State.Role(ctx.Guild.ID, id); err != nil {
				if _, err := services.GuildService.Role(id); err != nil {
					_, _ = ctx.Send(fmt.Sprintf("%s does not reference a valid role for this guild.", id))
					continue
				}
			}
			_, err := services.GuildService.CreateOrUpdateRole(geeksbot.Role{
				ID:       id,
				RoleType: "normal",
				Guild:    guild,
			})
			if err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("There was a problem adding <@&%s>", id), err)
				continue
			}
			added[id] = true
			count++
			_, _ = ctx.Send(fmt.Sprintf("Set <@&%s> as a normal role.", id))
		}
		_, _ = ctx.Send(fmt.Sprintf("Set %d roles to normal.", count))
	} else {
		_, _ = ctx.Send("Please include at least one role to remove from the moderator or admin lists.")
	}
}
