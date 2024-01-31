package catalogserver

import (
	"github.com/ConnectEverything/natster/internal/models"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
)

const (
	synadiaCloudServer = "tls://connect.ngs.global"
)

func generateConnectionFromOpts(opts *models.Options) (*nats.Conn, error) {
	ctxOpts := []natscontext.Option{
		natscontext.WithServerURL(synadiaCloudServer),
		natscontext.WithCreds(opts.Creds),
	}

	natsContext, err := natscontext.New("natster", false, ctxOpts...)

	if err != nil {
		return nil, err
	}

	conn, err := natsContext.Connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
