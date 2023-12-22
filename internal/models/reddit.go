package models

type GuildReddit struct {
	UserID  string `json:"userid"`
	GuildID string `json:"guildid"`
	Karma   int    `json:"value"`
}
