package inits

import (
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/webserver"
	"github.com/zekurio/kikuri/internal/util/static"
)

func InitWebserver(ctn di.Container) (err error) {

	cfg := ctn.Get(static.DiConfig).(models.Config)

	if cfg.Webserver.Enabled {
		ws, err := webserver.New(ctn)
		if err != nil {
			log.Error("Failed to initialize webserver", err)
			return err
		}

		go func() {
			if err = ws.ListenAndServeBlocking(); err != nil {
				log.Fatal("Failed starting up web server")
			}
		}()
		log.Info("Webserver started",
			"bindAddr", cfg.Webserver.Addr, "publicAddr", cfg.Webserver.PublicAddr)
	}

	return err
}
