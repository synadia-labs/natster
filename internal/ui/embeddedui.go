package natsterui

import (
	"embed"
	"io/fs"
	log "log/slog"
)

//go:embed web/dist
var app embed.FS

func GetEmbeddedUI() (fs.FS, error) {
	dist, err := fs.Sub(app, "web/dist")
	if err != nil {
		log.Error(
			"Embedded filesystem error",
			log.String("error", err.Error()),
		)
		return nil, err
	}
	return dist, nil
}
