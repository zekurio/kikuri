package inits

import (
	"strings"

	"github.com/zekurio/kikuri/internal/models"

	"github.com/charmbracelet/log"
	redis_pkg "github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/database/postgres"
	"github.com/zekurio/kikuri/internal/services/database/redis"
	"github.com/zekurio/kikuri/internal/util/static"
)

func InitDatabase(ctn di.Container) (db database.Database, err error) {
	cfg := ctn.Get(static.DiConfig).(models.Config)

	driver := strings.ToLower(cfg.Database.Type)

	switch driver {
	case "postgres":
		db, err = postgres.NewPostgres(cfg.Database.Postgres)
	default:
		log.Fatal("Invalid database driver specified")
	}

	if cfg.Cache.CacheDatabase {
		rd := ctn.Get(static.DiRedis).(*redis_pkg.Client)
		db = redis.NewRedisMiddleware(db, rd)
		log.Info("Database caching enabled")
	}

	return db, err
}
