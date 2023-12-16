package wsutil

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/zekurio/kikuri/internal/embedded"
	"github.com/zekurio/kikuri/internal/services/database/dberr"

	"github.com/charmbracelet/log"
)

func GetFS() (f http.FileSystem, err error) {
	fsys, err := fs.Sub(embedded.FrontendFiles, "webdist")
	if err != nil {
		return
	}
	_, err = fsys.Open("index.html")
	if os.IsNotExist(err) {
		log.Info("Using web files from web/dist")
		f = http.Dir("web/dist")
		err = nil
		return
	}
	if err != nil {
		return
	}
	log.Info("Using embedded web files")
	f = http.FS(fsys)
	return
}

func IsErrInternalOrNotFound(err error) error {
	if dberr.IsErrNotFound(err) {
		return fiber.ErrNotFound
	}

	return err
}

func GetQueryInt(ctx *fiber.Ctx, key string, def, min, max int) (int, error) {
	valStr := ctx.Query(key)
	if valStr == "" {
		return def, nil
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if val < min || (max > 0 && val > max) {
		return 0, fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("value of '%s' must be in bounds [%d, %d]", key, min, max))
	}

	return val, nil
}
