package auth

import "github.com/zekurio/daemon/pkg/arrayutils"

type AuthOrigin string

const (
	AuthOriginDiscord = AuthOrigin("origin:discord")
)

func (t Claims) IsAuthOrigin(origin AuthOrigin) bool {
	return arrayutils.Contains(t.Scopes, string(origin))
}
