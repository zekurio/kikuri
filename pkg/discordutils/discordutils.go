package discordutils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zekurio/daemon/internal/util/static"

	"github.com/bwmarrin/discordgo"
)

// GetInviteLink returns invite link for the app
func GetInviteLink(session *discordgo.Session) string {
	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&scope=%s&permissions=%d",
		session.State.User.ID, static.OAuthScopes, static.InvitePermission)
}

// GetDiscordSnowflakeCreationTime returns time when snowflake was created
func GetDiscordSnowflakeCreationTime(snowflake string) (time.Time, error) {
	sfI, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	timestamp := (sfI >> 22) + 1420070400000
	return time.Unix(timestamp/1000, timestamp), nil
}

// SendMessageDM sends message to user
func SendMessageDM(session *discordgo.Session, userID, message string) (msg *discordgo.Message, err error) {
	ch, err := session.UserChannelCreate(userID)
	if err != nil {
		return
	}
	msg, err = session.ChannelMessageSend(ch.ID, message)
	return
}

// SendEmbedMessageDM sends embed message to user
func SendEmbedMessageDM(session *discordgo.Session, userID string, embed *discordgo.MessageEmbed) (msg *discordgo.Message, err error) {
	ch, err := session.UserChannelCreate(userID)
	if err != nil {
		return
	}
	msg, err = session.ChannelMessageSendEmbed(ch.ID, embed)
	return
}

// IsAdmin checks if the given member has the admin permission set.
func IsAdmin(guild *discordgo.Guild, member *discordgo.Member) bool {
	if member == nil || guild == nil {
		return false
	}

	for _, r := range guild.Roles {
		if r.Permissions&0x8 != 0 {
			for _, mrID := range member.Roles {
				if r.ID == mrID {
					return true
				}
			}
		}
	}

	return false
}

// IsAFKChannel checks if the given channel is the guilds AFK channel.
func IsAFKChannel(guild *discordgo.Guild, channel *discordgo.Channel) bool {
	return guild.AfkChannelID == channel.ID
}

// GetMessages returns messages from channel
func GetMessages(session *discordgo.Session, channelID string, limit int) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message

	for {
		ms, err := session.ChannelMessages(channelID, limit, "", "", "")
		if err != nil {
			return nil, err
		}

		messages = append(messages, ms...)

		if len(ms) < limit {
			break
		}
	}

	return messages, nil
}

// GetMessageLink returns a message link
func GetMessageLink(msg *discordgo.Message, guildID string) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, msg.ChannelID, msg.ID)
}
