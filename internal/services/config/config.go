package config

import (
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/cryptoutils"
)

var DefaultConfig = Config{
	Discord: Discord{
		Token:            "",
		OwnerID:          "",
		GuildLimit:       -1,
		DisabledCommands: []string{},
	},
	Permissions: Permission{
		UserRules:  static.DefaultUserRules,
		AdminRules: static.DefaultAdminRules,
	},
	Webserver: Webserver{
		Enabled:    true,
		Addr:       ":80",
		PublicAddr: "http://localhost:80",
		DebugAddr:  "http://localhost:8081",
		AccessToken: AccessToken{
			Secret:          cryptoutils.ForceRandBase64Str(64),
			LifetimeSeconds: 10 * 60,
		},
		TLS: TLS{
			Enabled: false,
			Cert:    "",
			Key:     "",
		},
	},
}

type Discord struct {
	ClientID         string
	ClientSecret     string
	Token            string
	OwnerID          string
	GuildLimit       int
	DisabledCommands []string
}

type DatabaseCreds struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type CacheRedis struct {
	Addr     string
	Password string
	Type     int
}

type DatabaseType struct {
	Type     string
	Postgres DatabaseCreds
	Redis    CacheRedis
}

type Cache struct {
	Redis         CacheRedis
	CacheDatabase bool
	Lifetimes     CacheLifetimes
}

type CacheLifetimes struct {
	General,
	Guild,
	Member,
	User,
	Role,
	Channel,
	Emoji,
	Message,
	VoiceState,
	Presence string
}

type Permission struct {
	UserRules  []string
	AdminRules []string
}

type Webserver struct {
	Enabled     bool
	Addr        string
	PublicAddr  string
	DebugAddr   string
	AccessToken AccessToken
	TLS         TLS
}

type AccessToken struct {
	Secret          string
	LifetimeSeconds int
}

type TLS struct {
	Enabled bool
	Cert    string
	Key     string
}

type Config struct {
	Discord     Discord
	Database    DatabaseType
	Cache       Cache
	Permissions Permission
	Webserver   Webserver
}
