package acctelemetry

import "bytes"

type RealtimeUpdate struct {
}

func updateRealtimeUpdate(payload *bytes.Buffer, toUpdate *RealtimeUpdate) *RealtimeUpdate {
	return &RealtimeUpdate{}
}
