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

	_, _ = srv.nc.Subscribe(
		"*.natster.global.my.shares",
		handleMyShares(srv))

	_, _ = srv.nc.Subscribe(
		"*.natster.global.otc.generate",
		handleOtcGenerate(srv))

	_, _ = srv.nc.Subscribe(
		"*.natster.global.otc.claim",
		handleOtcClaim(srv))

	_, _ = srv.nc.Subscribe(
		"*.natster.global.whoami",
		handleWhoAmi(srv))

	_, _ = srv.nc.Subscribe(
		"*.natster.global.context.get",
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
