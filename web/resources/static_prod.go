//go:build !dev

package resources

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/benbjohnson/hashfs"
)

var (
	//go:embed static
	staticFiles embed.FS
	StaticSys   *hashfs.FS
)

func init() {
	// Strip the "static" prefix from embedded paths
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	StaticSys = hashfs.NewFS(fsys)
}

func Handler() http.Handler {
	slog.Debug("static assets are embedded from web/static/")
	return hashfs.FileServer(StaticSys)
}

func StaticPath(path string) string {
	return "/" + StaticSys.HashName(path)
}
