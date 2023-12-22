package models

import (
	"errors"

	"github.com/bwmarrin/snowflake"
	"github.com/zekurio/kikuri/pkg/hashutils"
)

type RedditAction string

const (
	RedditActionToggleRole  RedditAction = "TOGGLE_ROLE"
	RedditActionKick        RedditAction = "KICK"
	RedditActionBan         RedditAction = "BAN"
	RedditActionSendMessage RedditAction = "SEND_MESSAGE"
)

func (a RedditAction) Validate() bool {
	switch a {
	case RedditActionToggleRole, RedditActionKick, RedditActionBan, RedditActionSendMessage:
		return true
	default:
		return false
	}
}

type RedditTriggerType int

const (
	RedditTriggerBelow RedditTriggerType = iota
	RedditTriggerAbove

	RedditTriggerMax
)

func (tt RedditTriggerType) Validate() bool {
	return tt >= 0 && tt < RedditTriggerMax
}

type RedditRule struct {
	ID       snowflake.ID      `json:"id"`
	GuildID  string            `json:"guildid"`
	Trigger  RedditTriggerType `json:"trigger"`
	Value    int               `json:"value"`
	Action   RedditAction      `json:"action"`
	Argument string            `json:"argument"`
	Checksum string            `json:"-"`
}

func (r *RedditRule) Validate() error {
	if !r.Trigger.Validate() {
		return errors.New("invalid value for trigger")
	}
	if !r.Action.Validate() {
		return errors.New("invalid value for action")
	}

	return nil
}

func (r *RedditRule) CalculateChecksum() string {
	cop := *r
	cop.ID = 0
	r.Checksum = hashutils.Must(hashutils.SumMD5(&cop))
	return r.Checksum
}
