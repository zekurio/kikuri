package util

import (
	"fmt"
	"strings"

	"github.com/zekrotja/ken"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/util/static"
)

var allPermissions []string

func GetAllPermissions(cmdHandler *ken.Ken) []string {
	if allPermissions != nil {
		return allPermissions
	}

	cmds := cmdHandler.GetCommandInfo()

	// Create a copy and wrap it into a set
	perms := append([]string{}, static.AdditinalPerms...)

	for _, cmd := range cmds {
		rDomain := cmd.Implementations["Domain"]
		if len(rDomain) != 1 {
			continue
		}
		domain, ok := rDomain[0].(string)
		if !ok {
			continue
		}
		perms = append(perms, domain)

		rSubs := cmd.Implementations["SubDomains"]
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
				comb = fmt.Sprintf("%s.%s", domain, sub.Perm)
			}
			perms = append(perms, comb)
		}
	}

	allPermissions = perms
	return perms
}
