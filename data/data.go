package data

// CurrentVersion is the disklavier client build this server serves at
// /latest-client. A running client polls /latest-client-version and
// offers to update itself when the two differ, so bumping this is what
// tells deployed clients there is something new to fetch.
const CurrentVersion = "3.1.0"

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
