package models

import "github.com/zekurio/daemon/internal/util/static"

var DefaultConfig = Config{
	Discord: Discord{
		Token:            "",
		OwnerID:          "",
		GuildLimit:       -1,
		DisabledCommands: []string{},
	},
	Postgres: Postgres{
		Host: "localhost",
		Port: 5432,
	},
	Permissions: Permission{
		UserRules:  static.DefaultUserRules,
		AdminRules: static.DefaultAdminRules,
	},
	Webserver: Webserver{
		Enabled:    true,
		Addr:       ":8080",
		PublicAddr: "http://localhost:8080",
		TLS: WebserverTLS{
			Enabled: false,
			Cert:    "",
			Key:     "",
		},
	},
	Privacy: Privacy{
		NoticeURL: "",
		Contact: []Contact{
			{
				Title: "Example",
				Value: "Example Value",
				URL:   "https://example.com",
			},
		},
	},
}

type Discord struct {
	Token            string
	OwnerID          string
	GuildLimit       int
	DisabledCommands []string
}

type Postgres struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type Permission struct {
	UserRules  []string
	AdminRules []string
}

type Webserver struct {
	Enabled    bool
	Addr       string
	PublicAddr string
	TLS        WebserverTLS
}

type WebserverTLS struct {
	Enabled bool
	Cert    string
	Key     string
}

type Privacy struct {
	NoticeURL string
	Contact   []Contact
}

type Contact struct {
	Title string
	Value string
	URL   string
}

type Config struct {
	Discord     Discord
	Postgres    Postgres
	Permissions Permission
	Webserver   Webserver
	Privacy     Privacy
}
