//go:build with_tailscale

package main

import (
	"crypto/tls"
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

type TailscaleOptions struct {
	tailscalePort     int
	tailscaleNodeName string
	tailscaleHTTPS    bool
	tailscaleLogging  bool
	tailscaleAuthKey  string
}

var (
	TailscaleOpts = &TailscaleOptions{}
)

func init() {
	serveTailscale := natsterServer.Command("tailscale", "Serve webapp via tailscale")
	serveTailscale.Flag("addr", "Address to listen on").Default("8080").IntVar(&TailscaleOpts.tailscalePort)
	serveTailscale.Flag("https", "Tailscale will pull cert and serve on :443. Overrides port flag").Default("false").UnNegatableBoolVar(&TailscaleOpts.tailscaleHTTPS)
	serveTailscale.Flag("logs", "Output tailscale logs to stdout").Default("false").UnNegatableBoolVar(&TailscaleOpts.tailscaleLogging)
	serveTailscale.Flag("name", "Tailscale node name").Default("natster-ui").StringVar(&TailscaleOpts.tailscaleNodeName)
	serveTailscale.Flag("key", "Tailscale auth key").PlaceHolder("ts-auth-...").Envar("TS_AUTHKEY").Required().StringVar(&TailscaleOpts.tailscaleAuthKey)
	serveTailscale.Action(func(_ *fisk.ParseContext) error {
		return RunTailscale()
	})
}

func RunTailscale() error {
	dir, err := os.MkdirTemp(os.TempDir(), "natster-tailscale_*")
	if err != nil {
		return err
	}
	tempDirPath := filepath.Join(os.TempDir(), dir)

	s := &tsnet.Server{
		Dir:      tempDirPath,
		Hostname: TailscaleOpts.tailscaleNodeName,
		AuthKey:  TailscaleOpts.tailscaleAuthKey,
	}
	if !TailscaleOpts.tailscaleLogging {
		s.Logf = logger.Discard
	}

	defer s.Close()
	defer os.RemoveAll(tempDirPath)

	ln, err := s.Listen("tcp", fmt.Sprintf(":%d", TailscaleOpts.tailscalePort))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	if TailscaleOpts.tailscaleHTTPS || TailscaleOpts.tailscalePort == 443 {
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
