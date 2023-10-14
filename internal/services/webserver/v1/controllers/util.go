package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/webserver/v1/models"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/discordutils"
)

type UtilController struct {
	session *discordgo.Session
	cfg     config.Config
	ken     *ken.Ken
	st      *dgrs.State
}

func (c *UtilController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = ctn.Get(static.DiConfig).(config.Config)
	c.ken = ctn.Get(static.DiCommandHandler).(*ken.Ken)
	c.st = ctn.Get(static.DiState).(*dgrs.State)

	router.Get("/landingpageinfo", c.getLandingPageInfo)
	router.Get("/slashcommands", c.getSlashCommands)
}

// @Summary Landing Page Info
// @Description Returns general information for the landing page like the local invite parameters.
// @Tags Utilities
// @Accept json
// @Produce json
// @Success 200 {object} models.LandingPageResponse
// @Router /util/landingpageinfo [get]
func (c *UtilController) getLandingPageInfo(ctx *fiber.Ctx) error {
	res := new(models.LandingPageResponse)

	publicInvites := c.cfg.Webserver.LandingPage.ShowPublicInvites
	localInvite := c.cfg.Webserver.LandingPage.ShowLocalInvite

	if publicInvites {
		res.PublicCanaryInvite = static.PublicCanaryInvite
	}

	if localInvite {
		self, err := c.st.SelfUser()
		if err != nil {
			return err
		}
		res.LocalInvite = discordutils.GetInviteLink(self.ID)
	}

	return ctx.JSON(res)
}

// @Summary Slash Command List
// @Description Returns a list of registered slash commands and their description.
// @Tags Utilities
// @Accept json
// @Produce json
// @Success 200 {array} models.SlashCommandInfo "Wrapped in models.ListResponse"
// @Router /util/slashcommands [get]
func (c *UtilController) getSlashCommands(ctx *fiber.Ctx) error {
	cmdInfo := c.ken.GetCommandInfo()
	res := make([]*models.SlashCommandInfo, len(cmdInfo))

	for i, ci := range cmdInfo {
		res[i] = models.GetSlashCommandInfoFromCommand(ci)
	}

	return ctx.JSON(models.NewListResponse(res))
}
