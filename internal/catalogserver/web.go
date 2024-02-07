package catalogserver

import (
	"fmt"
	"io"
	log "log/slog"
	"net/http"

	natsterui "github.com/synadia-labs/natster/internal/ui"
)

func (hub *CatalogServer) startWebServer(port int) error {

	sockshub := natsterui.NewHub()
	go sockshub.Run()

	embedded, err := natsterui.GetEmbeddedUI()
	if err != nil {
		log.Error(
			"Failed to load embedded FS for UI",
			log.String("error", err.Error()),
		)
		return err
	}
	// TODO: use a router for the API functions
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		serveApi(hub, w, r)
	})
	http.Handle("/", http.FileServer(http.FS(embedded)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		natsterui.ServeWs(sockshub, w, r)
	})

	log.Info(
		"Starting HTTP Server",
		log.Int("port", port),
	)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Error(
		"Web Server Terminated",
		log.String("error", err.Error()),
	)

	return err
}

func serveApi(hub *CatalogServer, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Fill this in please!\n")
}
