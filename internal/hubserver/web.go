package hubserver

import (
	"fmt"
	"io"
	log "log/slog"
	"net/http"

	natsterui "github.com/ConnectEverything/natster/internal/ui"
)

func (hub *HubServer) startWebServer(port int) error {

	sockshub := natsterui.NewHub()
	go sockshub.Run()

	embedded, err := natsterui.GetEmbeddedUI()
	if err != nil {
		log.Error(
			"Failed to load embedded FS for UI",
			"error", err,
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
		"port", port,
	)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Error(
		"Web Server Terminated",
		"error", err,
	)

	return err
}

func serveApi(hub *HubServer, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Fill this in please!\n")
}
