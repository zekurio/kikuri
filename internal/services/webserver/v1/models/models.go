package models

import "github.com/zekurio/daemon/pkg/perms"

var Ok = &Status{Code: 200}

type Status struct {
	Code    int
	Message string
}

type GuildSettings struct {
	AutoRoles []string                    `json:"auto_roles"`
	AutoVoice []string                    `json:"auto_voice"`
	Perms     map[string]perms.PermsArray `json:"perms"`
}
