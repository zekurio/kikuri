package models

type GuildAPISettings struct {
	Enabled        bool   `json:"enabled"`
	AllowedOrigins string `json:"allowed_origins"`
	Protected      bool   `json:"protected"`
	TokenHash      string `json:"token_hash,omitempty"`
}
