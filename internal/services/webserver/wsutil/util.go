package wsutil

import (
	"github.com/charmbracelet/log"
	"github.com/zekurio/daemon/internal/util/embedded"
	"io/fs"
	"net/http"
	"os"
)

func GetFS() (f http.FileSystem, err error) {
	fsys, err := fs.Sub(embedded.FrontendFiles, "webdist")
	if err != nil {
		return
	}
	_, err = fsys.Open("index.html")
	if os.IsNotExist(err) {
		log.Info("Using web files from web/dist/web")
		f = http.Dir("web/dist/web")
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
