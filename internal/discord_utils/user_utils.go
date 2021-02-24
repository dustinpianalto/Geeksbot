package discord_utils

import "github.com/dustinpianalto/disgoman"

func GetDisplayName(ctx disgoman.Context, id string) (string, error) {
	member, err := ctx.Session.GuildMember(ctx.Guild.ID, id)
	if err != nil {
		return "", err
	}
	if member.Nick != "" {
		return member.Nick, nil
	}
	return member.User.Username, nil
}
