package exts

import (
	"github.com/dustinpianalto/geeksbot"
	"github.com/dustinpianalto/geeksbot/internal/exts/utils"
)

func AddCommandHandlers(g *geeksbot.Geeksbot) {
	// Arguments:
	// name - command name - string
	// desc - command description - string
	// owneronly - only allow owners to run - bool
	// hidden - hide command from non-owners - bool
	// perms - permissisions required - anpan.Permission (int)
	// type - command type, sets where the command is available
	// run - function to run - func(anpan.Context, []string) / CommandRunFunc
	_ = g.AddCommand(utils.UserCommand)
	_ = g.AddCommand(utils.SayCommand)
	_ = g.AddCommand(utils.GitCommand)
	_ = g.AddCommand(utils.InviteCommand)
	_ = g.AddCommand(utils.PingCommand)
}
