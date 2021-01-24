package discord_utils

import "github.com/dustinpianalto/disgoman"

func SendErrorMessage(ctx disgoman.Context, msg string, err error) {
	ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
		Context: ctx,
		Message: msg,
		Error:   err,
	}
}
