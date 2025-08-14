package lib

type (
	HttpMethod string
	RepeatMode string
)

const (
	GET  HttpMethod = "GET"
	PUT  HttpMethod = "PUT"
	POST HttpMethod = "POST"

	RepeatTrack   RepeatMode = "track"   // repeat the current track.
	RepeatContext RepeatMode = "context" // repeat the current context.
	RepeatOff     RepeatMode = "off"     // repeat off.
)
