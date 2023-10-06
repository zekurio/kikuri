package config

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
	ClientID         string
	ClientSecret     string
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

type Redis struct {
	Host      string
	Port      int
	Password  string
	Lifetimes CacheLifetimes
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
	APITokenKey string
	AccessToken AccessToken
	TLS         WebserverTLS
}

type AccessToken struct {
	Secret          string
	LifetimeSeconds int
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
	Redis       Redis
	Permissions Permission
	Webserver   Webserver
	Privacy     Privacy
}
