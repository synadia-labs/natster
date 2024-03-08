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

	binDir  string
	logging bool
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
	serve.Flag("binDir", "Directory holding natster binaries").Default("/tmp/natster_binaries/").StringVar(&Opts.binDir)
	serve.Flag("withLogging", "Logs web request").Default("false").UnNegatableBoolVar(&Opts.logging)
	serve.Action(func(_ *fisk.ParseContext) error {
		return RunServer()
	})

	natsterServer.MustParseWithUsage(os.Args[1:])
}

func RunServer() error {
	muxer := http.NewServeMux()

	info, err := os.Stat(Opts.binDir)
	if err != nil {
		fmt.Println("WARN no bin dir, not serving")
	} else {
		if info.IsDir() {
			binServe := http.FileServer(http.Dir(Opts.binDir))
			muxer.Handle("/dl/", http.StripPrefix("/dl", binServe))
		} else {
			fmt.Println("WARN user did not provide directory, not serving")
		}
	}

	var staticFS = fs.FS(static)
	htmlContent, err := fs.Sub(staticFS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	muxer.Handle("/", http.FileServer(http.FS(htmlContent)))

	fmt.Printf("Server started %s:%d\n", Opts.serveHost, Opts.servePort)
	if Opts.logging {
		return http.ListenAndServe(fmt.Sprintf("%s:%d", Opts.serveHost, Opts.servePort), logz(muxer))
	} else {
		return http.ListenAndServe(fmt.Sprintf("%s:%d", Opts.serveHost, Opts.servePort), muxer)
	}
}

func logz(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Host: %s | Path: %s | Query: %s\n", r.Host, r.URL.Path, r.URL.RawQuery)
		handler.ServeHTTP(w, r)
	})
}
