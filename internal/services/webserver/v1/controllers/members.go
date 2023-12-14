package controllers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	sharedmodels "github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/database/dberr"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/services/webserver/v1/models"
	"github.com/zekurio/kikuri/internal/services/webserver/wsutil"
	"github.com/zekurio/kikuri/internal/util"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/discordutils"
)

type GuildMembersController struct {
	session    *discordgo.Session
	db         database.Database
	st         *dgrs.State
	cmdHandler *ken.Ken
	pmw        *permissions.Permissions
	cfg        sharedmodels.Config
}

func (c *GuildMembersController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = ctn.Get(static.DiConfig).(sharedmodels.Config)
	c.db = ctn.Get(static.DiDatabase).(database.Database)
	c.pmw = ctn.Get(static.DiPermissions).(*permissions.Permissions)
	c.cmdHandler = ctn.Get(static.DiCommandHandler).(*ken.Ken)
	c.st = ctn.Get(static.DiState).(*dgrs.State)

	router.Get("/members", c.getMembers)
	router.Get("/:memberid", c.getMember)
	router.Get("/:memberid/permissions", c.getMemberPermissions)
	router.Get("/:memberid/permissions/allowed", c.getMemberPermissionsAllowed)
}

// @Summary Get Guild Member List
// @Description Returns a list of guild members.
// @Tags Members
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param after query string false "Request members after the given member ID."
// @Param limit query int false "The amount of results returned." default(100) minimum(1) maximum(1000)
// @Success 200 {array} models.Member "Wraped in models.ListResponse"
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/members [get]
func (c *GuildMembersController) getMembers(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	memb, _ := c.session.GuildMember(guildID, uid)
	if memb == nil {
		return fiber.ErrNotFound
	}

	after := ""
	limit := 0

	after = ctx.Query("after")
	limit, err = wsutil.GetQueryInt(ctx, "limit", 100, 1, 1000)
	if err != nil {
		return err
	}

	members, err := c.st.Members(guildID)
	if err != nil {
		return err
	}

	if filter := ctx.Query("filter"); filter != "" {
		filter = strings.ToLower(filter)
		var filteredMembers []*discordgo.Member
		for _, member := range members {
			if strings.Contains(strings.ToLower(member.Nick), filter) ||
				strings.Contains(strings.ToLower(member.User.Username), filter) ||
				strings.Contains(strings.ToLower(member.User.ID), filter) {
				filteredMembers = append(filteredMembers, member)
			}
		}
		members = filteredMembers
	} else if after != "" {
		for i := 0; i < len(members); i++ {
			if members[i].User.ID == after {
				members = members[i+1:]
				break
			}
		}
	}

	if limit > 0 && limit < len(members) {
		members = members[:limit]
	}

	fhmembers := make([]*models.Member, len(members))

	for i, m := range members {
		fhmembers[i] = models.MemberFromDiscordMember(m)
	}

	return ctx.JSON(models.NewListResponse(fhmembers))
}

// @Summary Get Guild Member
// @Description Returns a single guild member by ID.
// @Tags Members
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the member."
// @Success 200 {object} models.Member
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/{memberid} [get]
func (c *GuildMembersController) getMember(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	var memb *discordgo.Member

	if memb, _ = c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	guild, err := c.st.Guild(guildID, true)
	if err != nil {
		return err
	}

	memb, _ = c.session.GuildMember(guildID, memberID)
	if memb == nil {
		return fiber.ErrNotFound
	}

	memb.GuildID = guildID

	mm := models.MemberFromDiscordMember(memb)

	switch {
	case discordutils.IsAdmin(guild, memb):
		mm.Privilege = 1
	case guild.OwnerID == memberID:
		mm.Privilege = 2
	case c.cfg.Discord.OwnerID == memb.User.ID:
		mm.Privilege = 3
	}

	return ctx.JSON(mm)
}

// @Summary Get Guild Member Permissions
// @Description Returns the permission array of the given user.
// @Tags Members
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the member."
// @Success 200 {object} models.PermissionsResponse
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/{memberid}/permissions [get]
func (c *GuildMembersController) getMemberPermissions(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	if memb, _ := c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	perm, _, err := c.pmw.GetPerms(guildID, memberID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(&models.PermissionsResponse{
		Permissions: perm,
	})
}

// @Summary Get Guild Member Allowed Permissions
// @Description Returns all detailed permission DNS which the member is alloed to perform.
// @Tags Members
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the member."
// @Success 200 {array} string "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/{memberid}/permissions/allowed [get]
func (c *GuildMembersController) getMemberPermissionsAllowed(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	perms, _, err := c.pmw.GetPerms(guildID, memberID)
	if dberr.IsErrNotFound(err) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return err
	}

	all := util.GetAllPermissions(c.cmdHandler)
	allowed := make([]string, 0)
	for _, v := range all {
		if perms.Has(v) {
			allowed = append(allowed, v)
		}
	}

	return ctx.JSON(models.NewListResponse(allowed))
}
