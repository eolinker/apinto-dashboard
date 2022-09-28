package apinto_dashboard

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed builds/static
var staticDir embed.FS

func getStaticFiles() http.FileSystem {
	files, err := fs.Sub(staticDir, "builds/static")
	if err != nil {
		panic(err)
	}
	return http.FS(files)
}
