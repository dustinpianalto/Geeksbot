package exts

import (
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/geeksbot/internal/exts/guild"
	"github.com/dustinpianalto/geeksbot/internal/exts/requests"
	"github.com/dustinpianalto/geeksbot/internal/exts/utils"
)

func AddCommandHandlers(g *disgoman.CommandManager) {
	// Arguments:
	// name - command name - string
	// desc - command description - string
	// owneronly - only allow owners to run - bool
	// hidden - hide command from non-owners - bool
	// perms - permissisions required - anpan.Permission (int)
	// type - command type, sets where the command is available
	// run - function to run - func(anpan.Context, []string) / CommandRunFunc
	_ = g.AddCommand(utils.UserCommand)
	_ = g.AddCommand(utils.AddUserCommand)
	_ = g.AddCommand(utils.SayCommand)
	_ = g.AddCommand(utils.GitCommand)
	_ = g.AddCommand(utils.InviteCommand)
	_ = g.AddCommand(utils.PingCommand)
	_ = g.AddCommand(guild.AddPrefixCommand)
	_ = g.AddCommand(guild.RemovePrefixCommand)
	_ = g.AddCommand(guild.AddModeratorRoleCommand)
	_ = g.AddCommand(requests.RequestCommand)
}
