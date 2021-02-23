package discord_utils

import (
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/pkg/services"
)

func IsGuildMod(ctx disgoman.Context, user geeksbot.User) bool {
	discordCloser, err := ctx.Session.GuildMember(ctx.Guild.ID, user.ID)
	if err != nil {
		return false
	}
	guildRoles, err := services.GuildService.GuildRoles(geeksbot.Guild{ID: ctx.Guild.ID})
	if err != nil {
		return false
	}
	for _, role := range guildRoles {
		if role.RoleType == "moderator" {
			for _, rID := range discordCloser.Roles {
				if rID == role.ID {
					return true
				}
			}
		}
	}
	return false
}

func IsGuildAdmin(ctx disgoman.Context, user geeksbot.User) bool {
	discordCloser, err := ctx.Session.GuildMember(ctx.Guild.ID, user.ID)
	if err != nil {
		return false
	}
	guildRoles, err := services.GuildService.GuildRoles(geeksbot.Guild{ID: ctx.Guild.ID})
	if err != nil {
		return false
	}
	for _, role := range guildRoles {
		if role.RoleType == "admin" {
			for _, rID := range discordCloser.Roles {
				if rID == role.ID {
					return true
				}
			}
		}
	}
	return false
}
