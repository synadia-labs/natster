package hubserver

import (
	"strings"

	"github.com/ConnectEverything/natster/internal/models"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
)

func generateConnectionFromOpts(opts *models.Options) (*nats.Conn, error) {
	if len(strings.TrimSpace(opts.Servers)) == 0 {
		opts.Servers = "0.0.0.0:4222"
	}
	ctxOpts := []natscontext.Option{
		natscontext.WithServerURL(opts.Servers),
		natscontext.WithCreds(opts.Creds),
		natscontext.WithNKey(opts.Nkey),
		natscontext.WithCertificate(opts.TlsCert),
		natscontext.WithKey(opts.TlsKey),
		natscontext.WithCA(opts.TlsCA),
	}

	if opts.TlsFirst {
		ctxOpts = append(ctxOpts, natscontext.WithTLSHandshakeFirst())
	}

	if opts.Username != "" && opts.Password == "" {
		ctxOpts = append(ctxOpts, natscontext.WithToken(opts.Username))
	} else {
		ctxOpts = append(ctxOpts, natscontext.WithUser(opts.Username), natscontext.WithPassword(opts.Password))
	}

	natsContext, err := natscontext.New("natster", false, ctxOpts...)

	if err != nil {
		return nil, err
	}

	natsOpts, err := natsContext.NATSOptions()
	if err != nil {
		return nil, err
	}

	conn, err := nats.Connect(opts.Servers, natsOpts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
