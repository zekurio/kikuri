package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/ken"

	"github.com/zekurio/daemon/internal/inits"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/permissions"
	"github.com/zekurio/daemon/internal/util/embedded"
	"github.com/zekurio/daemon/internal/util/static"
)

var (
	flagConfigPath = flag.String("c", "config.toml", "Path to config file")
)

func main() {

	flag.Parse()

	if embedded.Release == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	diBuilder, err := di.NewBuilder()
	if err != nil {
		log.With(err).Fatal("Failed to create DI builder")
	}

	// Config
	err = diBuilder.Add(di.Def{
		Name: static.DiConfig,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.Parse(*flagConfigPath, "DAEMON_", config.DefaultConfig)
		},
	})
	if err != nil {
		log.With(err).Fatal("Config parsing failed")
	}

	// Database
	err = diBuilder.Add(di.Def{
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
	if err != nil && err.Error() == "unknown database driver" {
		log.With(err).Fatal("Database creation failed, unknown driver")
	} else if err != nil {
		log.With(err).Fatal("Database creation failed")
	}

	// Initialize discord bot session and shutdown routine
	err = diBuilder.Add(di.Def{
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
	if err != nil {
		return
	}

	// Permissions
	err = diBuilder.Add(di.Def{
		Name: static.DiPermissions,
		Build: func(ctn di.Container) (interface{}, error) {
			return permissions.InitPermissions(ctn), nil
		},
	})
	if err != nil {
		log.With(err).Fatal("Permissions creation failed")
	}

	// Ken
	err = diBuilder.Add(di.Def{
		Name: static.DiCommandHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitKen(ctn)
		},
		Close: func(obj interface{}) error {
			return obj.(*ken.Ken).Unregister()
		},
	})
	if err != nil {
		log.With(err).Fatal("Command handler creation failed")
	}

	// Scheduler
	err = diBuilder.Add(di.Def{
		Name: static.DiScheduler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitScheduler(ctn), nil
		},
	})
	if err != nil {
		log.With(err).Fatal("Scheduler creation failed")
	}

	// Webserver
	err = diBuilder.Add(di.Def{
		Name: static.DiWebserver,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitWebserver(ctn), nil
		},
	})
	if err != nil {
		log.With(err).Fatal("Webserver creation failed")
	}

	// Build dependency injection container
	ctn := diBuilder.Build()
	// Tear down dependency instances
	defer func(ctn di.Container) {
		err := ctn.DeleteWithSubContainers()
		if err != nil {
			log.With(err).Fatal("Failed to tear down dependency instances")
		}
	}(ctn)

	ctn.Get(static.DiCommandHandler)

	err = inits.InitDiscord(ctn)
	if err != nil {
		log.With(err).Fatal("Failed to initialize discord session")
	}

	ctn.Get(static.DiWebserver)

	ctn.Get(static.DiDatabase)

	// Block main go routine until one of the following
	// specified exit sys calls occure.
	log.Info("Started event loop. Stop with CTRL-C...")

	log.Info("Initialization finished")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
