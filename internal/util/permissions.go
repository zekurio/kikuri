package util

import (
	"fmt"
	"strings"

	"github.com/zekrotja/ken"
	"github.com/zekrotja/sop"
	"github.com/zekurio/daemon/internal/services/permissions"
	"github.com/zekurio/daemon/internal/util/static"
)

var allPerms sop.Enumerable[string]

func GetAllPermissions(cmdHandler *ken.Ken) sop.Enumerable[string] {
	if allPerms != nil {
		return allPerms
	}

	cmds := cmdHandler.GetCommandInfo()

	// Create a copy and wrap it into a set
	perms := sop.Set(append([]string{}, static.AdditionalRules...))

	for _, cmd := range cmds {
		rPerms := cmd.Implementations["Perm"]
		if len(rPerms) != 1 {
			continue
		}
		Perm, ok := rPerms[0].(string)
		if !ok {
			continue
		}
		perms.Push(Perm)

		rSubs := cmd.Implementations["SubPerms"]
		if len(rSubs) != 1 {
			continue
		}
		subs, ok := rSubs[0].([]permissions.SubCommandPerms)
		if !ok {
			continue
		}
		for _, sub := range subs {
			var comb string
			if strings.HasPrefix(sub.Perm, "/") {
				comb = sub.Perm[1:]
			} else {
				comb = fmt.Sprintf("%s.%s", Perm, sub.Perm)
			}
			perms.Push(comb)
		}
	}

	allPerms = perms
	return perms
}
