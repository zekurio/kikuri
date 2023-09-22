package inits

import (
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/webserver"
	"github.com/zekurio/daemon/internal/util/static"
)

func InitWebserver(ctn di.Container) (err error) {

	cfg := ctn.Get(static.DiConfig).(config.Config)

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
