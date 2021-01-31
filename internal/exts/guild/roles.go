package guild

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/discord_utils"
	"github.com/dustinpianalto/geeksbot/internal/services"
	"github.com/dustinpianalto/geeksbot/internal/utils"
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
		_, _ = ctx.Send(fmt.Sprintf("Added %d moderator %s.", count, utils.PluralizeString("role", count)))
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
		_, _ = ctx.Send(fmt.Sprintf("Added %d admin %s.", count, utils.PluralizeString("role", count)))
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
				if r, err := services.GuildService.Role(id); err != nil {
					_, _ = ctx.Send(fmt.Sprintf("%s does not reference a valid role for this guild.", id))
					continue
				} else {
					err = services.GuildService.DeleteRole(r)
					if err != nil {
						discord_utils.SendErrorMessage(ctx, "Something went wrong deleting the role", err)
					}
					_, _ = ctx.Send(fmt.Sprintf("Deleted <@&%s> as a no longer a valid role.", id))
					continue
				}
			}
			_, err := services.GuildService.CreateOrUpdateRole(geeksbot.Role{
				ID:       id,
				RoleType: "normal",
				Guild:    guild,
			})
			if err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("There was a problem updating <@&%s>", id), err)
				continue
			}
			added[id] = true
			count++
			_, _ = ctx.Send(fmt.Sprintf("Set <@&%s> as a normal role.", id))
		}
		_, _ = ctx.Send(fmt.Sprintf("Set %d %s to normal.", count, utils.PluralizeString("role", count)))
	} else {
		_, _ = ctx.Send("Please include at least one role to remove from the moderator or admin lists.")
	}
}

var MakeRoleSelfAssignableCommand = &disgoman.Command{
	Name:                "make-role-self-assignable",
	Aliases:             []string{"makesar", "addsar"},
	Description:         "Makes the passed in role self assignable by anyone",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              makeRoleSelfAssignableCommandFunc,
}

func makeRoleSelfAssignableCommandFunc(ctx disgoman.Context, args []string) {
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

			var role *discordgo.Role
			var err error
			if role, err = ctx.Session.State.Role(ctx.Guild.ID, id); err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%s does not reference a valid role for this guild", id), err)
				return
			}

			_, err = services.GuildService.CreateOrUpdateRole(geeksbot.Role{
				ID:       role.ID,
				RoleType: "sar",
				Guild:    guild,
			})
			if err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("There was a problem updating <@&%s>", role.ID), err)
				return
			}
			_, _ = ctx.Send(fmt.Sprintf("%s is now self assignable", role.Name))

		}
	} else {
		_, _ = ctx.Send("Please include at least one role to make self assignable")
	}

}

var RemoveSelfAssignableCommand = &disgoman.Command{
	Name:                "remove-self-assignable-role",
	Aliases:             []string{"removesar"},
	Description:         "Makes a role that was previously self assignable not so",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              removeSelfAssignableRoleCommandFunc,
}

func removeSelfAssignableRoleCommandFunc(ctx disgoman.Context, args []string) {
	removed := make(map[string]bool)
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
			if _, ok := removed[id]; ok {
				continue
			}

			var err error
			var role *discordgo.Role
			if role, err = ctx.Session.State.Role(ctx.Guild.ID, id); err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%s does not reference a valid role for this guild", id), err)
				return
			}
			_, err = services.GuildService.CreateOrUpdateRole(geeksbot.Role{
				ID:       role.ID,
				RoleType: "normal",
				Guild:    guild,
			})
			if err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("There was a problem updating <@&%s>", role.ID), err)
				return
			}
			_, _ = ctx.Send(fmt.Sprintf("%s's self assignability has been removed.", role.Name))
		}
	} else {
		_, _ = ctx.Send("Please include at least one role to make self assignable")
	}

}

var SelfAssignRoleCommand = &disgoman.Command{
	Name:                "giverole",
	Aliases:             []string{"iwant", "givetome"},
	Description:         "Assigns a person the passed in role if it is self assignable",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              selfAssignRoleCommandFunc,
}

func selfAssignRoleCommandFunc(ctx disgoman.Context, args []string) {
	added := make(map[string]bool)
	roles := append(args, ctx.Message.MentionRoles...)
	if len(roles) > 0 {
		for _, id := range roles {
			if strings.HasPrefix(id, "<@&") && strings.HasSuffix(id, ">") {
				continue
			}
			if _, ok := added[id]; ok {
				continue
			}
			var role *discordgo.Role
			var err error
			if role, err = ctx.Session.State.Role(ctx.Guild.ID, id); err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%s does not reference a valid role for this guild", id), err)
				return
			}
			if memberHasRole(ctx.Member, role.ID) {
				_, _ = ctx.Send(fmt.Sprintf("You already have the %s role silly...", role.Name))
				return
			}
			r, err := services.GuildService.Role(role.ID)
			if err != nil || r.RoleType != "sar" {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("You aren't allowed to assign yourself the %s role", role.Name), err)
				return
			}
			err = ctx.Session.GuildMemberRoleAdd(ctx.Guild.ID, ctx.User.ID, role.ID)
			if err != nil {
				discord_utils.SendErrorMessage(ctx, "There was a problem adding that role to you.", err)
				return
			}
			_, _ = ctx.Send(fmt.Sprintf("Congratulations! The %s role has been added to your... Ummm... Thing.", role.Name))
		}
	} else {
		_, _ = ctx.Send("Please include at least one role to make self assignable")
	}

}

var UnAssignRoleCommand = &disgoman.Command{
	Name:                "removerole",
	Aliases:             []string{"idon'twant"},
	Description:         "Removes a role from a person if the role is self assignable",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              unAssignRoleCommandFunc,
}

func unAssignRoleCommandFunc(ctx disgoman.Context, args []string) {
	removed := make(map[string]bool)
	roles := append(args, ctx.Message.MentionRoles...)
	if len(roles) > 0 {
		for _, id := range roles {
			if strings.HasPrefix(id, "<@&") && strings.HasSuffix(id, ">") {
				continue
			}
			if _, ok := removed[id]; ok {
				continue
			}

			var role *discordgo.Role
			var err error
			if role, err = ctx.Session.State.Role(ctx.Guild.ID, id); err != nil {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("%s does not reference a valid role for this guild", id), err)
				return
			}
			if !memberHasRole(ctx.Member, role.ID) {
				_, _ = ctx.Send(fmt.Sprintf("I can't remove the %s role from you because you don't have it...", role.Name))
				return
			}
			r, err := services.GuildService.Role(role.ID)
			if err != nil || r.RoleType != "sar" {
				discord_utils.SendErrorMessage(ctx, fmt.Sprintf("You aren't allowed to assign yourself the %s role", role.Name), err)
				return
			}
			err = ctx.Session.GuildMemberRoleRemove(ctx.Guild.ID, ctx.User.ID, role.ID)
			if err != nil {
				discord_utils.SendErrorMessage(ctx, "There was a problem removing that role from your account", err)
				return
			}
			_, _ = ctx.Send(fmt.Sprintf("Sad to see you go... but the %s role has been removed.", role.Name))
		}
	} else {
		_, _ = ctx.Send("Please include at least one role to make self assignable")
	}

}

func memberHasRole(m *discordgo.Member, id string) bool {
	for _, r := range m.Roles {
		if r == id {
			return true
		}
	}
	return false
}
