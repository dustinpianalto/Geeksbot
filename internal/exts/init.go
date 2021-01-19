package exts

import (
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff/internal/exts/fun"
	"github.com/dustinpianalto/goff/internal/exts/guild_management"
	"github.com/dustinpianalto/goff/internal/exts/roles"
	"github.com/dustinpianalto/goff/internal/exts/tags"
	"github.com/dustinpianalto/goff/internal/exts/tasks"
	"github.com/dustinpianalto/goff/internal/exts/user_management"
	"github.com/dustinpianalto/goff/internal/exts/utils"

	"github.com/dustinpianalto/goff/internal/exts/p_interpreter"
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
	_ = h.AddCommand(tasks.AddReminderCommand)
	_ = h.AddCommand(tags.AddTagCommand)
	_ = h.AddCommand(tags.TagCommand)
	_ = h.AddCommand(roles.MakeRoleSelfAssignableCommand)
	_ = h.AddCommand(roles.RemoveSelfAssignableCommand)
	_ = h.AddCommand(roles.SelfAssignRoleCommand)
	_ = h.AddCommand(roles.UnAssignRoleCommand)
	_ = h.AddCommand(p_interpreter.PCommand)
	_ = h.AddCommand(fun.InterleaveCommand)
	_ = h.AddCommand(fun.DeinterleaveCommand)
	_ = h.AddCommand(fun.GenerateRPNCommand)
	_ = h.AddCommand(fun.ParseRPNCommand)
	_ = h.AddCommand(fun.SolveCommand)
	_ = h.AddCommand(user_management.KickUserCommand)
	_ = h.AddCommand(user_management.BanUserCommand)
	_ = h.AddCommand(user_management.UnbanUserCommand)
	_ = h.AddCommand(guild_management.SetLoggingChannelCommand)
	_ = h.AddCommand(guild_management.GetLoggingChannelCommand)
	_ = h.AddCommand(guild_management.SetWelcomeChannelCommand)
	_ = h.AddCommand(guild_management.GetWelcomeChannelCommand)
	_ = h.AddCommand(guild_management.AddGuildCommand)
	_ = h.AddCommand(guild_management.SetPuzzleChannelCommand)
	_ = h.AddCommand(guild_management.GetPuzzleChannelCommand)
	_ = h.AddCommand(guild_management.SetPuzzleRoleCommand)
	_ = h.AddCommand(guild_management.GetPuzzleRoleCommand)

}
