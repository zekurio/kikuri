package models

import (
	"time"

	"github.com/zekurio/kikuri/pkg/perms"
)

var Ok = &Status{Code: 200}

type Status struct {
	Code    int
	Message string
}

type GuildSettings struct {
	AutoRoles  []string               `json:"auto_roles"`
	AutoVoice  []string               `json:"auto_voice"`
	Perms      map[string]perms.Array `json:"perms"`
	APIEnabled bool                   `json:"api_enabled"`
}

type GuildSettingsEmpty struct {
	APIEnabled bool `json:"api_enabled"`
}

type AccessTokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}
