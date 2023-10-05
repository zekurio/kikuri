package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/permissions"
)

type GuildMembersController struct {
	session    *discordgo.Session
	cfg        config.Config
	db         database.Database
	pmw        *permissions.Permissions
	cmdHandler *ken.Ken
	st         *dgrs.State
}
