package webserver

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/models"
	v1 "github.com/zekurio/kikuri/internal/services/webserver/v1"
	"github.com/zekurio/kikuri/internal/services/webserver/v1/controllers"
	"github.com/zekurio/kikuri/internal/services/webserver/wsutil"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/debug"
)

type WebServer struct {
	app       *fiber.App
	cfg       models.Config
	container di.Container
}

func New(ctn di.Container) (ws *WebServer, err error) {

	ws = new(WebServer)

	ws.container = ctn

	ws.cfg = ctn.Get(static.DiConfig).(models.Config)

	ws.app = fiber.New(fiber.Config{
		AppName:               "kikuri",
		DisableStartupMessage: true,
		ProxyHeader:           "X-Forwarded-For",
	})

	if debug.Enabled() {
		const corsOrigin = "http://localhost:5173"
		log.Warnf("CORS enabled for address %s", corsOrigin)
		ws.app.Use(cors.New(cors.Config{
			AllowOrigins:     ws.cfg.Webserver.DebugAddr,
			AllowHeaders:     "authorization, content-type, set-cookie, cookie, server",
			AllowMethods:     "GET, POST, PUT, PATCH, POST, DELETE, OPTIONS",
			AllowCredentials: true,
		}))
	}

	new(controllers.InviteController).Setup(ws.container, ws.app.Group("/invite"))

	ws.registerRouter(new(v1.Router), []string{"/api/v1", "/api"})

	fs, err := wsutil.GetFS()
	if err != nil {
		return
	}

	ws.app.Use(filesystem.New(filesystem.Config{
		Root:         fs,
		Browse:       true,
		Index:        "index.html",
		MaxAge:       3600,
		NotFoundFile: "index.html",
	}))

	return
}

func (ws *WebServer) registerRouter(router *v1.Router, routes []string, middlewares ...fiber.Handler) {
	router.SetContainer(ws.container)
	for _, r := range routes {
		router.Route(ws.app.Group(r, middlewares...))
	}
}

func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.cfg.Webserver.TLS

	if tls.Enabled {
		return ws.app.ListenTLS(ws.cfg.Webserver.Addr, tls.Cert, tls.Key)
	}

	return ws.app.Listen(ws.cfg.Webserver.Addr)
}
