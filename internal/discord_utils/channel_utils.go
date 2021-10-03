package discord_utils

import "github.com/dustinpianalto/disgoman"

func GetChannelName(ctx disgoman.Context, id string) string {
	channel, err := ctx.Session.Channel(id)
	if err != nil {
		return ""
	}
	return channel.Name
}
