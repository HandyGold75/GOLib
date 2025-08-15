package lib

import "errors"

type (
	HttpMethod  string
	RepeatMode  string
	TopItemType string
	TimeRange   string
)

const (
	GET    HttpMethod = "GET"
	PUT    HttpMethod = "PUT"
	POST   HttpMethod = "POST"
	DELETE HttpMethod = "DELETE"
)

const (
	RepeatTrack   RepeatMode = "track"   // repeat the current track.
	RepeatContext RepeatMode = "context" // repeat the current context.
	RepeatOff     RepeatMode = "off"     // repeat off.

	TimeRangeLongTerm   TimeRange = "long_term"   // calculated from ~1 year of data and including all new data as it becomes available
	TimeRangeMediumTerm TimeRange = "medium_term" // approximately last 6 months.
	TimeRangeShortTerm  TimeRange = "short_term"  // approximately last 4 weeks.
)

var Errors = struct{ UnexpectedResponse error }{
	UnexpectedResponse: errors.New("unexpected response"),
}
