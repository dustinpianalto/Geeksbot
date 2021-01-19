package exts

import (
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot/internal/exts/utils"
)

func AddCommandHandlers(h *disgoman.CommandManager) {
	// Arguments:
	// name - command name - string
	// desc - command description - string
	// owneronly - only allow owners to run - bool
	// hidden - hide command from non-owners - bool
	// perms - permissisions required - anpan.Permission (int)
	// type - command type, sets where the command is available
	// run - function to run - func(anpan.Context, []string) / CommandRunFunc
	_ = h.AddCommand(utils.UserCommand)
	_ = h.AddCommand(utils.SayCommand)
	_ = h.AddCommand(utils.GitCommand)
	_ = h.AddCommand(utils.InviteCommand)
	_ = h.AddCommand(utils.PingCommand)
}
