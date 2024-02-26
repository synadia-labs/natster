package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/choria-io/fisk"
)

//go:embed dist
var static embed.FS

type Options struct {
	servePort int
	serveHost string
}

var (
	VERSION = "dev"
	COMMIT  string
	Opts    = &Options{}

	natsterServer = fisk.New("natster-server", "used to run natster ui")
)

func main() {
	natsterServer.Author("Synadia Communications")
	natsterServer.UsageWriter(os.Stdout)
	natsterServer.HelpFlag.Short('h')
	natsterServer.WithCheats().CheatCommand.Hidden()
	natsterServer.Version(fmt.Sprintf("v%s [%s]", VERSION, COMMIT))

	serve := natsterServer.Command("serve", "Serve webapp")
	serve.Flag("host", "Host to serve on").StringVar(&Opts.serveHost)
	serve.Flag("port", "Port to listen on").Default("8080").IntVar(&Opts.servePort)
	serve.Action(func(_ *fisk.ParseContext) error {
		return RunServer()
	})

	natsterServer.MustParseWithUsage(os.Args[1:])
}

func RunServer() error {
	var staticFS = fs.FS(static)
	htmlContent, err := fs.Sub(staticFS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(htmlContent))

	fmt.Printf("Server started %s:%d\n", Opts.serveHost, Opts.servePort)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", Opts.serveHost, Opts.servePort), fs)
}
