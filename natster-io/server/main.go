package main

import (
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/choria-io/fisk"
	"tailscale.com/tsnet"
	"tailscale.com/types/logger"
)

//go:embed dist
var static embed.FS

type Options struct {
	servePort int
	serveHost string

	tailscalePort     int
	tailscaleNodeName string
	tailscaleHTTPS    bool
	tailscaleLogging  bool
	tailscaleAuthKey  string
}

var (
	VERSION = "dev"
	COMMIT  string
	Opts    = &Options{}
)

func main() {
	natsterServer := fisk.New("natster-server", "used to run natster ui")
	natsterServer.Author("Synadia Communications")
	natsterServer.UsageWriter(os.Stdout)
	natsterServer.HelpFlag.Short('h')
	natsterServer.WithCheats().CheatCommand.Hidden()
	natsterServer.Version(fmt.Sprintf("v%s [%s]", VERSION, COMMIT))

	serve := natsterServer.Command("serve", "Serve webapp")
	serve.Flag("host", "Host to serve on").Default("").StringVar(&Opts.serveHost)
	serve.Flag("addr", "Address to listen on").Default("8080").IntVar(&Opts.servePort)
	serve.Action(func(_ *fisk.ParseContext) error {
		return RunServer()
	})

	serveTailscale := natsterServer.Command("tailscale", "Serve webapp via tailscale")
	serveTailscale.Flag("addr", "Address to listen on").Default("8080").IntVar(&Opts.tailscalePort)
	serveTailscale.Flag("https", "Tailscale will pull cert and serve on :443. Overrides port flag").Default("false").UnNegatableBoolVar(&Opts.tailscaleHTTPS)
	serveTailscale.Flag("logs", "Output tailscale logs to stdout").Default("false").UnNegatableBoolVar(&Opts.tailscaleLogging)
	serveTailscale.Flag("name", "Tailscale node name").Default("natster-ui").StringVar(&Opts.tailscaleNodeName)
	serveTailscale.Flag("key", "Tailscale auth key").PlaceHolder("ts-auth-...").Envar("TS_AUTHKEY").Required().StringVar(&Opts.tailscaleAuthKey)
	serveTailscale.Action(func(_ *fisk.ParseContext) error {
		return RunTailscale()
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
	return http.ListenAndServe(fmt.Sprintf(":%d", Opts.servePort), fs)
}

func RunTailscale() error {
	f, err := os.CreateTemp(os.TempDir(), "natster-tailscale_*")
	if err != nil {
		return err
	}
	tempDirPath := filepath.Join(os.TempDir(), f.Name())

	s := &tsnet.Server{
		Dir:      tempDirPath,
		Hostname: Opts.tailscaleNodeName,
		AuthKey:  Opts.tailscaleAuthKey,
	}
	if !Opts.tailscaleLogging {
		s.Logf = logger.Discard
	}

	defer s.Close()
	defer os.RemoveAll(tempDirPath)

	ln, err := s.Listen("tcp", fmt.Sprintf(":%d", Opts.tailscalePort))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	if Opts.tailscaleHTTPS || Opts.tailscalePort == 443 {
		ln = tls.NewListener(ln, &tls.Config{
			GetCertificate: lc.GetCertificate,
		})
	}

	var staticFS = fs.FS(static)
	htmlContent, err := fs.Sub(staticFS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(htmlContent))

	return http.Serve(ln, fs)
}
