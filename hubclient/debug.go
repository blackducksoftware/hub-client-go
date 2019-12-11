package hubclient

import (
	log "github.com/sirupsen/logrus"
)

type HubClientDebug uint16

const (
	HubClientDebugTimings HubClientDebug = 1 << iota
	HubClientDebugContent
)

func debugReportBytes(bodyBytes []byte, debugFlags HubClientDebug) {
	if debugFlags&HubClientDebugContent != 0 {
		log.Debugf("START DEBUG: --------------------------------------------------------------------------- \n\n")
		log.Debugf("TEXT OF RESPONSE: \n %s", string(bodyBytes[:]))
		log.Debugf("END DEBUG: --------------------------------------------------------------------------- \n\n\n\n")
	}
}
