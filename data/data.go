package data

import "time"

const (
	CurrentVersion = "3.1.0"
)

var (
	PerformanceDate = mustParseTime("2006-01-02 15:04:05 MST", "2024-03-23 1:00:00 CDT")
)

const (
	PathLatestClientVersion   = "GET /latest-client-version"
	PathLatestClientDownload  = "GET /latest-client"
	PathPlayerWS              = "GET /ws"
	PathControllerWS          = "GET /controller-ws"
	PathSchedulePerformance   = "POST /performances"
	PathFeaturedPerformances  = "GET /performances/featured"
	PathScheduledPerformances = "GET /performances/scheduled"
	PathRestart               = "POST /performances/{id}/restart"
	PathAdvance               = "POST /performances/{id}/advance"
	PathMIDIFile              = "GET /performances/{id}/midi/{filename}"
	PathBeginPerformance      = "POST /performances/{id}/begin"
	PathDeletePerformance     = "DELETE /performances/{id}"
)

func mustParseTime(fmt string, t string) time.Time {
	if out, err := time.Parse(fmt, t); err != nil {
		panic(err)
	} else {
		return out
	}
}
