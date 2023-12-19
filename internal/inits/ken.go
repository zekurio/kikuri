package inits

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/state"

	"github.com/zekurio/kikuri/internal/middlewares"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/slashcommands"
	"github.com/zekurio/kikuri/internal/usercommands"
	"github.com/zekurio/kikuri/internal/util/static"
)

func InitKen(ctn di.Container) (k *ken.Ken, err error) {

	s := ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	st := ctn.Get(static.DiState).(*dgrs.State)
	p := ctn.Get(static.DiPermissions).(*permissions.Permissions)

	k, err = ken.New(s, ken.Options{
		EmbedColors: ken.EmbedColors{
			Default: static.ColorDefault,
			Error:   static.ColorRed,
		},
		DependencyProvider: ctn,
		OnSystemError:      systemErrorHandler,
		OnCommandError:     commandErrorHandler,
		State:              state.NewDgrs(st),
	})

	if err != nil {
		return nil, err
	}

	// register slashcommands
	err = k.RegisterCommands(
		// slashcommands
		new(slashcommands.Profile),
		new(slashcommands.Autovoice),
		new(slashcommands.Guild),
		new(slashcommands.Perms),
		new(slashcommands.Vote),

		// usercommands
		new(usercommands.About),
	)
	if err != nil {
		return nil, err
	}

	err = k.RegisterMiddlewares(
		p,
		middlewares.NewDisableCommandsMiddleware(ctn),
	)

	return k, err
}

func systemErrorHandler(context string, err error, args ...interface{}) {
	log.Error("Ken system error")
}

func commandErrorHandler(err error, ctx *ken.Ctx) {
	ctx.Defer()

	if err == ken.ErrNotDMCapable {
		ctx.FollowUpError("This command can not be used in DMs.", "").Send()
		return
	}

	ctx.FollowUpError(
		fmt.Sprintf("The command execution failed unexpectedly:\n```\n%s\n```", err.Error()),
		"Command execution failed").Send()
}
