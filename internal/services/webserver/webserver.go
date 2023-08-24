package webserver

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/config"
	v1 "github.com/zekurio/daemon/internal/services/webserver/v1"
	"github.com/zekurio/daemon/internal/services/webserver/v1/controllers"
	"github.com/zekurio/daemon/internal/util/static"
)

type WebServer struct {
	app       *fiber.App
	cfg       config.Config
	container di.Container
}

func New(ctn di.Container) (ws *WebServer, err error) {
	ws = new(WebServer)

	ws.container = ctn

	ws.cfg = ctn.Get(static.DiConfig).(config.Config)

	ws.app = fiber.New(fiber.Config{
		AppName:               "daemon",
		DisableStartupMessage: true,
		ProxyHeader:           "X-Forwarded-For",
	})

	new(controllers.InviteController).Setup(ws.container, ws.app.Group("/invite"))
	ws.registerRouter(new(v1.Router), []string{"/api"})

	return ws, nil
}

func (ws *WebServer) registerRouter(router Router, routes []string, middlewares ...fiber.Handler) {
	router.SetContainer(ws.container)
	for _, r := range routes {
		router.Route(ws.app.Group(r, middlewares...))
	}
}

func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.cfg.Webserver.TLS

	if tls.Enabled {
		if tls.Cert == "" || tls.Key == "" {
			return errors.New("cert file and key file must be specified")
		}
		return ws.app.ListenTLS(ws.cfg.Webserver.Addr, tls.Cert, tls.Key)
	}

	return ws.app.Listen(ws.cfg.Webserver.Addr)
}
