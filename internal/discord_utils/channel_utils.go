package discord_utils

import "github.com/dustinpianalto/disgoman"

func GetChannelName(ctx disgoman.Context, id string) (string, error) {
	channel, err := ctx.Session.Channel(id)
	if err != nil {
		return "", err
	}
	return channel.Name, nil
}
