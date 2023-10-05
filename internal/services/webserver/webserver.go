package webserver

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/config"
	v1 "github.com/zekurio/daemon/internal/services/webserver/v1"
	"github.com/zekurio/daemon/internal/services/webserver/v1/controllers"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/debug"
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

	if debug.Enabled() {
		const corsOrigin = "http://localhost:5173"
		log.Warnf("CORS enabled for address %s", corsOrigin)
		ws.app.Use(cors.New())
	}

	new(controllers.InviteController).Setup(ws.container, ws.app.Group("/invite"))
	ws.registerRouter(new(v1.Router), []string{"/api"}) // moving to /api since we only have one version

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
		log.Infof("Starting webserver on %s (TLS enabled)", ws.cfg.Webserver.Addr)
		return ws.app.ListenTLS(ws.cfg.Webserver.Addr, tls.Cert, tls.Key)
	}

	log.Infof("Starting webserver on %s", ws.cfg.Webserver.Addr)
	return ws.app.Listen(ws.cfg.Webserver.Addr)
}
