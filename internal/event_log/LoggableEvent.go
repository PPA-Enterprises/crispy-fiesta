package event_log
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event = string
const(
	Created = "Created"
	Edited = "Edited"
	Deleted = "Deleted"
)

type Change struct {
	Old string
	New string
}

type LoggableEvent struct {
	ID primitive.ObjectID `json:"_id,omitemty" bson:"_id,omitemty"`
	EventType Event `json:"event_type" bson:"event_type"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	Editor string `json:"editor" bson:"editor"`
	Changes map[Field]Change `json:"changes" bson:"changes"`
}
