package wsutil

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/zekurio/daemon/internal/util/embedded"
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
