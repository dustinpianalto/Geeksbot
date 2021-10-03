package discord_utils

import "github.com/dustinpianalto/disgoman"

func GetDisplayName(ctx disgoman.Context, id string) string {
	member, err := ctx.Session.GuildMember(ctx.Guild.ID, id)
	if err != nil {
		return ""
	}
	if member.Nick != "" {
		return member.Nick
	}
	return member.User.Username
}
