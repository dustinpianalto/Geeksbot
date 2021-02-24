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
	_ = g.AddCommand(guild.AddAdminRoleCommand)
	_ = g.AddCommand(guild.RemoveModRoleCommand)
	_ = g.AddCommand(guild.MakeRoleSelfAssignableCommand)
	_ = g.AddCommand(guild.RemoveSelfAssignableCommand)
	_ = g.AddCommand(guild.SelfAssignRoleCommand)
	_ = g.AddCommand(guild.UnAssignRoleCommand)
	_ = g.AddCommand(requests.RequestCommand)
	_ = g.AddCommand(requests.CloseCommand)
	_ = g.AddCommand(requests.ListCommand)
	_ = g.AddCommand(requests.ViewCommand)
	_ = g.AddCommand(requests.CommentCommand)
}
