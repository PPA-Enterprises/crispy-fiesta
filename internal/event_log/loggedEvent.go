package event_log

type loggedEvent struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	EventType event `json:"event_type" bson:"event_type"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	Editor string `json:"editor" bson:"editor"`
	EditorID primitive.ObjectID `json:"editor_id" bson:"editor_id"`
	Changes map[field]change `json:"changes" bson:"changes"`
}
