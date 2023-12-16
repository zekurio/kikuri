package inits

import (
	"reflect"
	"time"

	"github.com/zekurio/kikuri/internal/models"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/timeutils"
)

func getLifetimes(cfg models.Config) (dgrs.Lifetimes, bool, error) {
	lifetimes := cfg.Cache.Lifetimes

	var target dgrs.Lifetimes

	vlt := reflect.ValueOf(lifetimes)
	vtg := reflect.ValueOf(&target)

	set := false

	for i := 0; i < vlt.NumField(); i++ {
		ds := vlt.Field(i).String()
		if ds == "" {
			continue
		}

		d, err := timeutils.ParseDuration(ds)
		if err != nil {
			return dgrs.Lifetimes{}, false, err
		}

		if d == 0 {
			continue
		}

		vtg.Elem().FieldByName(vlt.Type().Field(i).Name).Set(reflect.ValueOf(d))
		set = true
	}

	return target, set, nil
}

func InitState(container di.Container) (s *dgrs.State, err error) {
	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	rd := container.Get(static.DiRedis).(*redis.Client)
	cfg := container.Get(static.DiConfig).(models.Config)

	lf, set, err := getLifetimes(cfg)
	if err != nil {
		return nil, err
	}

	if !set {
		lf.General = 7 * 24 * time.Hour
		log.Warn("No cache lifetimes set in config. Using default of 7 days.")
	}

	lf.OverrrideZero = true

	return dgrs.New(dgrs.Options{
		RedisClient:    rd,
		DiscordSession: session,
		FetchAndStore:  true,
		Lifetimes:      lf,
	})
}
