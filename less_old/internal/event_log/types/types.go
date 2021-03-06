package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Change struct {
	Old interface{}
	New interface{}
}

type Editor struct {
	Oid primitive.ObjectID
	Name string
	Collection string
}

type NormalizedLoggedEvent struct {
	ID string `json:"_id" bson:"_id"`
	EventType string `json:"event_type" bson:"event_type"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	Editor string `json:"editor" bson:"editor"`
	EditorID string `json:"editor_id" bson:"editor_id"`
	Persisted bool `json:"persisted" bson:"persisted"`
	Changes map[string]Change `json:"changes" bson:"changes"`
}
