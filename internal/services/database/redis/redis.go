package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type RedisDatabase struct {
	redisClient *redis.Client
	pgDatabase  *sqlx.DB
}

func NewRedisDatabase(redisClient *redis.Client, pgDatabase *sqlx.DB) *RedisDatabase {
	return &RedisDatabase{
		redisClient: redisClient,
		pgDatabase:  pgDatabase,
	}
}

func (r *RedisDatabase) Close() error {
	if err := r.redisClient.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}

	if err := r.pgDatabase.Close(); err != nil {
		return fmt.Errorf("failed to close Postgres database: %w", err)
	}

	return nil
}

func (r *RedisDatabase) GetAutoRoles(guildID string) ([]string, error) {
	key := fmt.Sprintf("guild:%s:auto_roles", guildID)

	val, err := r.redisClient.Get(r.redisClient.Context(), key).Result()
	if err == redis.Nil {
		// Cache miss, fetch from Postgres
		var roles []string
		err = r.pgDatabase.Select(&roles, "SELECT role_id FROM auto_roles WHERE guild_id = $1", guildID)
		if err != nil {
			return nil, fmt.Errorf("failed to get auto roles from Postgres: %w", err)
		}

		// Cache the result in Redis
		if len(roles) > 0 {
			err = r.redisClient.Set(r.redisClient.Context(), key, roles, 5*time.Minute).Err()
			if err != nil {
				return nil, fmt.Errorf("failed to cache auto roles in Redis: %w", err)
			}
		}

		return roles, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get auto roles from Redis: %w", err)
	}

	// Cache hit, unmarshal the result
	var roles []string
	err = json.Unmarshal([]byte(val), &roles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal auto roles from Redis: %w", err)
	}

	return roles, nil
}

func (r *RedisDatabase) SetAutoRoles(guildID string, roleIDs []string) error {
	key := fmt.Sprintf("guild:%s:auto_roles", guildID)

	// Update Postgres
	_, err := r.pgDatabase.Exec("DELETE FROM auto_roles WHERE guild_id = $1", guildID)
	if err != nil {
		return fmt.Errorf("failed to delete auto roles from Postgres: %w", err)
	}

	if len(roleIDs) > 0 {
		query, args, err := sqlx.In("INSERT INTO auto_roles (guild_id, role_id) VALUES (?, ?)", guildID, roleIDs)
		if err != nil {
			return fmt.Errorf("failed to create SQL query for auto roles: %w", err)
		}

		_, err = r.pgDatabase.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to insert auto roles into Postgres: %w", err)
		}
	}

	// Update Redis
	if len(roleIDs) > 0 {
		err = r.redisClient.Set(r.redisClient.Context(), key, roleIDs, 5*time.Minute).Err()
		if err != nil {
			return fmt.Errorf("failed to cache auto roles in Redis: %w", err)
		}
	} else {
		err = r.redisClient.Del(r.redisClient.Context(), key).Err()
		if err != nil {
			return fmt.Errorf("failed to delete auto roles from Redis: %w", err)
		}
	}

	return nil
}
