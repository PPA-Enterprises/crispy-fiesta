package event_log
import (
	//"time"
	"internal/event_log/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type loggedEvent struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	EventType event `json:"event_type" bson:"event_type"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	Editor string `json:"editor" bson:"editor"`
	EditorID primitive.ObjectID `json:"editor_id" bson:"editor_id"`
	Persisted bool `json:"persisted" bson:"persisted"`
	Changes map[field]change `json:"changes" bson:"changes"`
}
func (self *loggedEvent) normalize() *types.NormalizedLoggedEvent {
	return &types.NormalizedLoggedEvent {
		ID: self.ID.Hex(),
		EventType: self.EventType,
		Timestamp: self.Timestamp,
		Editor: self.Editor,
		EditorID: self.EditorID.Hex(),
		changes: self.Changes,
	}
}
