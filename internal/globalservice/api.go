package globalservice

import (
	"strings"
)

func (srv *GlobalService) startApiSubscriptions() error {
	_, _ = srv.nc.Subscribe(
		"*.natster.global.events.put",
		handleEventPut(srv))

	_, _ = srv.nc.Subscribe(
		"*.natster.global.heartbeats.put",
		handleHeartbeat(srv))

	_, _ = srv.nc.Subscribe(
		"*.natster.global.stats",
		handleStats(srv))

	return nil
}

func extractAccountKey(subject string) string {
	tokens := strings.Split(subject, ".")
	if len(tokens) > 0 {
		return tokens[0]
	}
	return "???"
}
