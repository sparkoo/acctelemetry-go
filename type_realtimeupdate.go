package acctelemetry

import "bytes"

type RealtimeUpdate struct {
}

func createRealtimeUpdate(payload *bytes.Buffer) *RealtimeUpdate {
	return &RealtimeUpdate{}
}
