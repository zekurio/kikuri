package usercommands

import (
	"github.com/zekrotja/ken"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/slashcommands"
)

type About struct {
	slashcommands.Profile
}

var (
	_ ken.UserCommand          = (*About)(nil)
	_ permissions.CommandPerms = (*About)(nil)
)

func (c *About) TypeUser() {}

func (c *About) Name() string {
	return "about"
}
