package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/ken"

	"github.com/zekurio/daemon/internal/inits"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/permissions"
	"github.com/zekurio/daemon/internal/services/webserver/auth"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/debug"
)

var (
	flagConfigPath = flag.String("c", "config.toml", "Path to config file")
	flagDebug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {

	flag.Parse()

	debug.SetEnabled(*flagDebug)

	if debug.Enabled() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	diBuilder, err := di.NewBuilder()
	if err != nil {
		log.With("err", err).Fatal("Failed to create DI builder")
	}

	// Config
	diBuilder.Add(di.Def{
		Name: static.DiConfig,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.Parse(*flagConfigPath, "DAEMON_", config.DefaultConfig)
		},
	})

	// Redis
	diBuilder.Add(di.Def{
		Name: static.DiRedis,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(static.DiConfig).(config.Config)
			return redis.NewClient(&redis.Options{
				Addr:     cfg.Cache.Redis.Addr,
				Password: cfg.Cache.Redis.Password,
				DB:       cfg.Cache.Redis.Type,
			}), nil
		},
	})

	// Database
	diBuilder.Add(di.Def{
		Name: static.DiDatabase,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitDatabase(ctn)
		},
		Close: func(obj interface{}) error {
			d := obj.(database.Database)
			log.Info("Shutting down database connection...")
			err := d.Close()
			if err != nil {
				return err
			}
			return nil
		},
	})

	// Initialize discord bot session and shutdown routine
	diBuilder.Add(di.Def{
		Name: static.DiDiscordSession,
		Build: func(ctn di.Container) (interface{}, error) {
			return discordgo.New("")
		},
		Close: func(obj interface{}) error {
			session := obj.(*discordgo.Session)
			log.Info("Shutting down bot session...")
			err := session.Close()
			if err != nil {
				return err
			}
			return nil
		},
	})

	// Initialize Discord OAuth Module
	diBuilder.Add(di.Def{
		Name: static.DiDiscordOAuth,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitDiscordOAuth(ctn), nil
		},
	})

	// Initialize auth refresh token handler
	diBuilder.Add(di.Def{
		Name: static.DiAuthRefreshTokenHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewDBRefreshTokenHandler(ctn), nil
		},
	})

	// Initialize auth access token handler
	diBuilder.Add(di.Def{
		Name: static.DiAuthAccessTokenHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewJWTAccessTokenHandler(ctn), nil
		},
	})

	// Initialize auth API token handler
	diBuilder.Add(di.Def{
		Name: static.DiAuthAPITokenHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewDBAPITokenHandler(ctn), nil
		},
	})

	// Initialize OAuth API handler implementation
	diBuilder.Add(di.Def{
		Name: static.DiOAuthHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewRefreshTokenRequestHandler(ctn), nil
		},
	})

	// Initialize access token authorization middleware
	diBuilder.Add(di.Def{
		Name: static.DiAuthMiddleware,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewAccessTokenMiddleware(ctn), nil
		},
	})

	// Initialize State
	diBuilder.Add(di.Def{
		Name: static.DiState,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitState(ctn)
		},
	})

	// Permissions
	diBuilder.Add(di.Def{
		Name: static.DiPermissions,
		Build: func(ctn di.Container) (interface{}, error) {
			return permissions.InitPermissions(ctn), nil
		},
	})

	// Ken
	diBuilder.Add(di.Def{
		Name: static.DiCommandHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitKen(ctn)
		},
		Close: func(obj interface{}) error {
			return obj.(*ken.Ken).Unregister()
		},
	})

	// Scheduler
	diBuilder.Add(di.Def{
		Name: static.DiScheduler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitScheduler(ctn), nil
		},
	})

	// Webserver
	diBuilder.Add(di.Def{
		Name: static.DiWebserver,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitWebserver(ctn), nil
		},
	})

	// Build dependency injection container
	ctn := diBuilder.Build()
	// Tear down dependency instances
	defer func(ctn di.Container) {
		err := ctn.DeleteWithSubContainers()
		if err != nil {
			log.With("err", err).Fatal("Failed to tear down dependency instances")
		}
	}(ctn)

	ctn.Get(static.DiCommandHandler)

	err = inits.InitDiscord(ctn)
	if err != nil {
		log.With("err", err).Fatal("Failed to initialize discord session")
	}

	ctn.Get(static.DiWebserver)

	ctn.Get(static.DiDatabase)

	// Block main go routine until one of the following
	// specified exit sys calls occure.
	log.Info("Started event loop. Stop with CTRL-C...")

	log.Info("Initialization finished")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
