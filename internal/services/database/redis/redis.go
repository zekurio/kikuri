package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-redis/redis/v8"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/util"
)

const (
	keyGuildAutoRoles = "GUILD:AUTOROLES"
	keyGuildAutoVoice = "GUILD:AUTOVOICE"

	keyAPIToken = "API:TOKEN"
)

type RedisMiddleware struct {
	database.Database

	client *redis.Client
}

var _ database.Database = (*RedisMiddleware)(nil)

// NewRedisMiddleware creates a new RedisMiddleware, which wraps a Database
// and adds Redis caching to it.
func NewRedisMiddleware(db database.Database, rd *redis.Client) *RedisMiddleware {
	return &RedisMiddleware{
		Database: db,
		client:   rd,
	}
}

func (r *RedisMiddleware) Close() error {
	if err := r.client.Close(); err != nil {
		log.Error("failed to close Redis client: %w", err)
		return err
	}

	return nil
}

func (r *RedisMiddleware) GetGuildAutoRoles(guildID string) ([]string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildAutoRoles, guildID)

	valC, err := r.client.Get(context.Background(), key).Result()
	val := strings.Split(valC, ";")
	if err == redis.Nil {
		val, err = r.Database.GetGuildAutoRoles(guildID)
		if err != nil {
			return nil, err
		}

		err = r.client.Set(context.Background(), key, strings.Join(val, ";"), 0).Err()
		return val, err
	}
	if err != nil {
		return nil, err
	}

	if valC == "" {
		return []string{}, nil
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildAutoRoles(guildID string, roleIDs []string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildAutoRoles, guildID)

	err := r.Database.SetGuildAutoRoles(guildID, roleIDs)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), key, strings.Join(roleIDs, ";"), 0).Err()
}

func (r *RedisMiddleware) GetGuildAutoVoice(guildID string) ([]string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildAutoVoice, guildID)

	valC, err := r.client.Get(context.Background(), key).Result()
	val := strings.Split(valC, ";")
	if err == redis.Nil {
		val, err = r.Database.GetGuildAutoVoice(guildID)
		if err != nil {
			return nil, err
		}
	}

	val, err = r.Database.GetGuildAutoVoice(guildID)
	if err != nil {
		return nil, err
	}

	err = r.client.Set(context.Background(), key, strings.Join(val, ";"), 0).Err()
	return val, err
}

func (r *RedisMiddleware) SetGuildAutoVoice(guildID string, channelIDs []string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildAutoVoice, guildID)

	err := r.Database.SetGuildAutoVoice(guildID, channelIDs)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), key, strings.Join(channelIDs, ";"), 0).Err()
}

func (r *RedisMiddleware) SetAPIToken(token models.APITokenEntry) error {
	var key = fmt.Sprintf("%s:%s", keyAPIToken, token.UserID)

	data, err := util.Marshal(token)
	if err != nil {
		return err
	}

	if err := r.client.Set(context.Background(), key, data, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetAPIToken(token)
}

func (r *RedisMiddleware) GetAPIToken(userID string) (models.APITokenEntry, error) {
	var key = fmt.Sprintf("%s:%s", keyAPIToken, userID)

	res, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		token, err := r.Database.GetAPIToken(userID)
		if err != nil {
			return token, err
		}

		data, err := util.Marshal(token)
		if err != nil {
			return token, err
		}

		err = r.client.Set(context.Background(), key, data, 0).Err()
		return token, err
	}

	return util.Unmarshal[models.APITokenEntry](res)
}

func (r *RedisMiddleware) DeleteAPIToken(userID string) error {
	var key = fmt.Sprintf("%s:%s", keyAPIToken, userID)

	if err := r.client.Del(context.Background(), key).Err(); err != nil {
		return err
	}

	return r.Database.DeleteAPIToken(userID)
}
