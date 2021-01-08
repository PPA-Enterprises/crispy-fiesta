package types

import ("context")
type LoggableEvent interface {
	Log(ctx context.Context, collection string) *Deliverable
}

type Deliverable interface {
	normalize() *NormalizedLoggedEvent
}

type change struct {
	Old string
	New string
}

type NormalizedLoggedEvent struct {
	ID string `json:"_id" bson:"_id"`
	EventType string `json:"event_type" bson:"event_type"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	Editor string `json:"editor" bson:"editor"`
	EditorID string `json:"editor_id" bson:"editor_id"`
	Persisted bool `json:"persisted" bson:"persisted"`
	Changes map[string]change `json:"changes" bson:"changes"`
}
