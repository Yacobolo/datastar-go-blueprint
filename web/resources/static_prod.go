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
	// StaticSys is the hash-based filesystem for static assets.
	StaticSys *hashfs.FS
)

func init() {
	// Strip the "static" prefix from embedded paths
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	StaticSys = hashfs.NewFS(fsys)
}

// Handler returns an HTTP handler for serving static assets.
func Handler() http.Handler {
	slog.Debug("static assets are embedded from web/static/") //nolint:sloglint // Initialization code, global logger acceptable
	return hashfs.FileServer(StaticSys)
}

// StaticPath returns the hashed path for a static asset.
func StaticPath(path string) string {
	return "/" + StaticSys.HashName(path)
}
