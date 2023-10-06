package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/database/models"
)

const (
	keyGuildAutoRoles = "GUILD:AUTOROLES"
	keyGuildAutoVoice = "GUILD:AUTOVOICE"
	keyGuildAPI       = "GUILD:API"

	keyUserAPIToken = "USER:APITOKEN"
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
		return fmt.Errorf("failed to close Redis client: %w", err)
	}

	if err := r.Close(); err != nil {
		return fmt.Errorf("failed to close Postgres database: %w", err)
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

func (r *RedisMiddleware) GetGuildAPI(guildID string) (settings models.GuildAPISettings, err error) {
	var key = fmt.Sprintf("%s:%s", keyGuildAPI, guildID)

	resStr, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		if settings, err = r.Database.GetGuildAPI(guildID); err != nil {
			return
		}
		var resB []byte
		resB, err = json.Marshal(settings)
		if err != nil {
			return
		}
		if err = r.client.Set(context.Background(), key, resB, 0).Err(); err != nil {
			return
		}
		return
	}

	err = json.Unmarshal([]byte(resStr), &settings)

	return
}

func (r *RedisMiddleware) SetGuildAPI(guildID string, settings models.GuildAPISettings) error {
	var key = fmt.Sprintf("%s:%s", keyGuildAPI, guildID)

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	if err = r.client.Set(context.Background(), key, data, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildAPI(guildID, settings)
}

func (r *RedisMiddleware) GetAPIToken(userID string) (t models.APITokenEntry, err error) {
	var key = fmt.Sprintf("%s:%s", keyUserAPIToken, userID)

	resStr, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		if t, err = r.Database.GetAPIToken(userID); err != nil {
			return
		}
		var resB []byte
		resB, err = json.Marshal(t)
		if err != nil {
			return
		}
		if err = r.client.Set(context.Background(), key, resB, 0).Err(); err != nil {
			return
		}
		return
	}

	err = json.Unmarshal([]byte(resStr), &t)

	return
}

func (r *RedisMiddleware) SetAPIToken(token models.APITokenEntry) (err error) {
	var key = fmt.Sprintf("%s:%s", keyUserAPIToken, token.UserID)

	data, err := json.Marshal(token)
	if err != nil {
		return
	}

	if err = r.client.Set(context.Background(), key, data, 0).Err(); err != nil {
		return
	}

	return r.Database.SetAPIToken(token)
}

func (r *RedisMiddleware) DeleteAPIToken(userID string) (err error) {
	var key = fmt.Sprintf("%s:%s", keyUserAPIToken, userID)

	if err = r.client.Del(context.Background(), key).Err(); err != nil {
		return
	}

	return r.Database.DeleteAPIToken(userID)
}
