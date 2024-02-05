package globalservice

import (
	"strings"

	"github.com/nats-io/nats.go"
)

func (srv *GlobalService) startApiSubscriptions() error {
	_, _ = srv.nc.Subscribe(
		"*.natster.global.events.put",
		handleEventPut(srv))

	_, _ = srv.nc.Subscribe(
		"*.natster.global.heartbeats.put",
		handleHeartbeat(srv))

	return nil
}

func handleEventPut(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		// NOOP
	}
}

func extractAccountKey(subject string) string {
	tokens := strings.Split(subject, ".")
	if len(tokens) > 0 {
		return tokens[0]
	}
	return "???"
}
