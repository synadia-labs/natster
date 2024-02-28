package globalservice

import (
	"strings"
)

func (srv *GlobalService) startApiSubscriptions() error {
	_, _ = srv.nc.QueueSubscribe(
		"*.natster.global.events.put", "globalservice",
		handleEventPut(srv))

	_, _ = srv.nc.QueueSubscribe(
		"*.natster.global.heartbeats.put", "globalservice",
		handleHeartbeat(srv))

	_, _ = srv.nc.QueueSubscribe(
		"*.natster.global.stats", "globalservice",
		handleStats(srv))

	_, _ = srv.nc.QueueSubscribe(
		"*.natster.global.my.shares", "globalservice",
		handleMyShares(srv))

	_, _ = srv.nc.QueueSubscribe(
		"*.natster.global.otc.generate", "globalservice",
		handleOtcGenerate(srv))

	_, _ = srv.nc.QueueSubscribe(
		"*.natster.global.otc.claim", "globalservice",
		handleOtcClaim(srv))

	_, _ = srv.nc.QueueSubscribe(
		"*.natster.global.whoami", "globalservice",
		handleWhoAmi(srv))

	_, _ = srv.nc.QueueSubscribe(
		"*.natster.global.context.get", "globalservice",
		handleGetContext(srv))

	return nil
}

func extractAccountKey(subject string) string {
	tokens := strings.Split(subject, ".")
	if len(tokens) > 0 {
		return tokens[0]
	}
	return "???"
}
