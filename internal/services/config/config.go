package config

import "github.com/zekurio/daemon/internal/util/static"

var DefaultConfig = Config{
	Discord: DiscordConfig{
		Token:            "",
		OwnerID:          "",
		GuildLimit:       -1,
		DisabledCommands: []string{},
	},
	Postgres: PostgresConfig{
		Host: "localhost",
		Port: 5432,
	},
	Permissions: PermissionRules{
		UserRules:  static.DefaultUserRules,
		AdminRules: static.DefaultAdminRules,
	},
	Webserver: WebserverConfig{
		Enabled:    true,
		Addr:       ":8080",
		PublicAddr: "http://localhost:8080",
		TLS: TLSConfig{
			Enabled: false,
		},
	},
}

type DiscordConfig struct {
	Token            string
	OwnerID          string
	GuildLimit       int
	DisabledCommands []string
}

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type PermissionRules struct {
	UserRules  []string
	AdminRules []string
}

type WebserverConfig struct {
	Enabled    bool
	Addr       string
	PublicAddr string
	TLS        TLSConfig
}

type TLSConfig struct {
	Enabled bool
	Cert    string
	Key     string
}

type Config struct {
	Discord     DiscordConfig
	Postgres    PostgresConfig
	Permissions PermissionRules
	Webserver   WebserverConfig
}
