package permissions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"

	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/database/dberr"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/discordutils"
	"github.com/zekurio/daemon/pkg/perms"
	"github.com/zekurio/daemon/pkg/roleutils"
)

type Permissions struct {
	db  database.Database
	cfg config.Config
	s   *discordgo.Session
	st  *dgrs.State
}

var _ PermsProvider = (*Permissions)(nil)

func InitPermissions(ctn di.Container) *Permissions {
	return &Permissions{
		db:  ctn.Get(static.DiDatabase).(database.Database),
		cfg: ctn.Get(static.DiConfig).(config.Config),
		s:   ctn.Get(static.DiDiscordSession).(*discordgo.Session),
		st:  ctn.Get(static.DiState).(*dgrs.State),
	}
}

func (p *Permissions) Before(ctx *ken.Ctx) (next bool, err error) {
	cmd, ok := ctx.Command.(CommandPerms)
	if !ok {
		next = true
		return
	}

	if ctx.User() == nil {
		return
	}

	ok, _, err = p.HasPerms(ctx.GetSession(), ctx.GetEvent().GuildID, ctx.User().ID, cmd.Perm())

	if err != nil {
		return false, err
	}

	if !ok {
		err = ctx.RespondError("You are not permitted to use this command!", "Missing Permission")
		return
	}

	return true, err
}

func (p *Permissions) GetPerms(session *discordgo.Session, guildID, userID string) (perm perms.Array, override bool, err error) {

	if guildID != "" {
		perm, err = p.GetMemberPerms(session, guildID, userID)
		if err != nil && !dberr.IsErrNotFound(err) {
			return
		}
	} else {
		perm = make(perms.Array, 0)
	}

	if p.cfg.Discord.OwnerID == userID {
		perm = perms.Array{"+dm.*"}
		override = true
		return
	}

	if guildID != "" {
		guild, err := p.st.Guild(guildID)
		if err != nil {
			return perms.Array{}, false, err
		}

		member, err := p.st.Member(guildID, userID)
		if err != nil {
			return perms.Array{}, false, err
		}

		if userID == guild.OwnerID || (member != nil && discordutils.IsAdmin(guild, member)) {
			var defaultAdminPerms []string
			defaultAdminPerms = p.cfg.Permissions.AdminRules
			if defaultAdminPerms == nil {
				defaultAdminPerms = static.DefaultAdminRules
			}

			perm = perm.Merge(defaultAdminPerms, false)

			override = true

		}
	}

	var defaultUserPerms []string
	defaultUserPerms = p.cfg.Permissions.UserRules
	if defaultUserPerms == nil {
		defaultUserPerms = static.DefaultUserRules
	}

	perm = perm.Merge(defaultUserPerms, false)

	return perm, override, nil

}

func (p *Permissions) GetMemberPerms(session *discordgo.Session, guildID string, memberID string) (perms.Array, error) {
	guildPerms, err := p.db.GetPermissions(guildID)
	if err != nil {
		return nil, err
	}
	membRoles, err := roleutils.GetSortedMemberRoles(session, guildID, memberID, false, true)
	if err != nil {
		return nil, err
	}

	var res perms.Array
	for _, r := range membRoles {
		if p, ok := guildPerms[r.ID]; ok {
			if res == nil {
				res = p
			} else {
				res = res.Merge(p, true)
			}
		}
	}

	return res, nil
}

func (p *Permissions) HasPerms(session *discordgo.Session, guildID, userID, pm string) (ok, override bool, err error) {
	perm, override, err := p.GetPerms(session, guildID, userID)
	if err != nil {
		return false, false, err
	}

	return perm.Has(pm), override, nil
}

func (p *Permissions) HasSubCmdPerms(ctx ken.Context, subPM string, explicit bool, message ...string) (ok bool, err error) {

	cmd, cok := ctx.GetCommand().(CommandPerms)
	if !cok {
		return
	}

	var pm string
	if strings.HasPrefix(subPM, "/") {
		pm = subPM[1:]
	} else {
		pm = cmd.Perm() + "." + subPM
	}

	if explicit {
		pm = "!" + pm
	}

	msg := "Sorry, you are not permitted to use this command!"

	if len(message) != 0 {
		msg = message[0]
	}

	permOk, override, err := p.HasPerms(ctx.GetSession(), ctx.GetEvent().GuildID, ctx.User().ID, pm)
	if err != nil {
		return false, err
	}

	if !permOk && (explicit && !override) {
		err = ctx.FollowUpError(msg, "").Send().Error
		return
	}

	return true, nil

}
