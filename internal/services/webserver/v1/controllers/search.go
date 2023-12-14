package controllers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekurio/kikuri/internal/services/webserver/v1/models"
	"github.com/zekurio/kikuri/internal/services/webserver/wsutil"
	"github.com/zekurio/kikuri/internal/util/static"
)

type SearchController struct {
	session *discordgo.Session
	st      *dgrs.State
}

func (c *SearchController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.st = ctn.Get(static.DiState).(*dgrs.State)

	router.Get("", c.getSearch)
}

// @Summary Global Search
// @Description Search through guilds and members by ID, name or displayname.
// @Tags Search
// @Accept json
// @Produce json
// @Param query query string true "The search query (either ID, name or displayname)."
// @Param limit query int false "The maximum amount of result items (per group)." default(50) minimum(1) maximum(100)
// @Success 200 {object} models.SearchResult
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /search [get]
func (c *SearchController) getSearch(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)
	query := strings.ToLower(ctx.Query("query"))
	limit, err := wsutil.GetQueryInt(ctx, "limit", 50, 1, 100)
	if err != nil {
		return
	}

	if query == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query must be set")
	}

	result := &models.SearchResult{
		Guilds:  make([]*models.GuildReduced, 0),
		Members: make([]*models.Member, 0),
	}

	guilds, err := c.st.Guilds()
	if err != nil {
		return
	}

	for _, g := range guilds {
		// check if user is member of guild
		_, err := c.st.Member(g.ID, uid)
		if err != nil {
			continue
		}

		currentResults := len(result.Guilds) + len(result.Members)
		if currentResults >= limit {
			break
		}

		if strings.Contains(strings.ToLower(g.ID), query) || strings.Contains(strings.ToLower(g.Name), query) {
			// check for limit reached

			result.Guilds = append(result.Guilds, &models.GuildReduced{
				ID:   g.ID,
				Name: g.Name,
			})
		} else {
			members, err := c.st.Members(g.ID)
			if err != nil {
				return err
			}

			for _, m := range members {
				if strings.Contains(strings.ToLower(m.User.ID), query) || strings.Contains(strings.ToLower(m.User.Username), query) || strings.Contains(strings.ToLower(m.Nick), query) {
					result.Members = append(result.Members, models.MemberFromDiscordMember(m))
				}
			}
		}
	}

	return ctx.JSON(result)
}
