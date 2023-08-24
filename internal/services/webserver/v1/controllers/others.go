package controllers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/webserver/v1/models"
	"github.com/zekurio/daemon/internal/util"
	"github.com/zekurio/daemon/internal/util/embedded"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/discordutils"
	"runtime"
	"time"
)

type OthersController struct {
	session *discordgo.Session
}

func (c *OthersController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)

	router.Get("/sysinfo", c.getSysinfo)
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
