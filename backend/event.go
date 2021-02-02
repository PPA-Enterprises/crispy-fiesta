package PPA
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event = string
type Field = string

const(
	Created Event = "Created"
	Edited Event = "Edited"
	Deleted Event = "Deleted"
)
type Change struct {
	Old interface{} `json:"old" bson:"old"`
	New interface{} `json:"new" bson:"new"`
}

type Editor struct {
	OID primitive.ObjectID `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Collection string `json:"collection" bson:"collection"`
}

type LogEvent struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	EventType Event `json:"event_type" bson:"event_type"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Editor string `json:"editor" bson:"editor"`
	EditorID primitive.ObjectID `json:"editor_id" bson:"editor_id"`
	Persisted bool `json:"persisted" bson:"persisted"`
	Changes map[Field] Change `json:"changes" bson:"changes"`
}


