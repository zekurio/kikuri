package controllers

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/webserver/auth"
	"github.com/zekurio/daemon/internal/services/webserver/v1/models"
	"github.com/zekurio/daemon/internal/util"
	"github.com/zekurio/daemon/internal/util/embedded"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/discordutils"
)

type OthersController struct {
	cfg     config.Config
	st      *dgrs.State
	db      database.Database
	ken     *ken.Ken
	authMw  auth.Middleware
	session *discordgo.Session
}

func (c *OthersController) Setup(ctn di.Container, router fiber.Router) {
	c.cfg = ctn.Get(static.DiConfig).(config.Config)
	c.st = ctn.Get(static.DiState).(*dgrs.State)
	c.authMw = ctn.Get(static.DiAuthMiddleware).(auth.Middleware)
	c.db = ctn.Get(static.DiDatabase).(database.Database)
	c.ken = ctn.Get(static.DiCommandHandler).(*ken.Ken)
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)

	router.Get("/me", c.authMw.Handle, c.getMe)
	router.Get("/sysinfo", c.getSysinfo)
	router.Get("/privacyinfo", c.getPrivacyInfo)
}

// @Summary Me
// @Description Returns the user object of the currently authenticated user.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {object} apiModels.User
// @Router /me [get]
func (c *OthersController) getMe(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	user, err := c.st.User(uid)
	if err != nil {
		return err
	}

	created, _ := discordutils.GetDiscordSnowflakeCreationTime(user.ID)

	res := &models.User{
		User:      user,
		AvatarURL: user.AvatarURL(""),
		CreatedAt: created,
		BotOwner:  uid == c.cfg.Discord.OwnerID,
	}

	return ctx.JSON(res)
}

// @Summary System Information
// @Description Returns general global system information.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {object} models.SystemInfo
// @Router /sysinfo [get]
func (c *OthersController) getSysinfo(ctx *fiber.Ctx) error {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := int64(time.Since(util.StatsStartupTime).Seconds())

	guilds := c.session.State.Guilds

	res := &models.SystemInfo{
		Version:    embedded.AppVersion,
		CommitHash: embedded.AppCommit,
		GoVersion:  runtime.Version(),

		Uptime:    uptime,
		UptimeStr: fmt.Sprintf("%d", uptime),

		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		CPUs:        runtime.NumCPU(),
		GoRoutines:  runtime.NumGoroutine(),
		StackUse:    memStats.StackInuse,
		StackUseStr: fmt.Sprintf("%d", memStats.StackInuse),
		HeapUse:     memStats.HeapInuse,
		HeapUseStr:  fmt.Sprintf("%d", memStats.HeapInuse),

		BotUserID: c.session.State.User.ID,
		BotInvite: discordutils.GetInviteLink(c.session),

		Guilds: len(guilds),
	}

	return ctx.JSON(res)
}

// @Summary Privacy Information
// @Description Returns information about the privacy policy.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {object} models.Privacy
// @Router /privacyinfo [get]
func (c *OthersController) getPrivacyInfo(ctx *fiber.Ctx) error {
	return ctx.JSON(c.cfg.Privacy)
}

// @Summary All Permissions
// @Description Return a list of all available permissions.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {array} string "Wrapped in models.ListResponse"
// @Router /allpermissions [get]
func (c *OthersController) getAllPermissions(ctx *fiber.Ctx) error {
	all := util.GetAllPermissions(c.ken)
	return ctx.JSON(models.NewListResponse(all.Unwrap()))
}
