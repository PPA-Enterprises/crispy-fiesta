package event_log
import (
	//"time"
	"context"
	"internal/event_log/types"
	"internal/db"
	"internal/uid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type event = string
const(
	created = "Created"
	edited = "Edited"
	deleted = "Deleted"
	failed = "Failed to Log"
)

type change struct {
	Old string
	New string
}

type logEvent struct {
	EventType event `json:"event_type" bson:"event_type"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	Editor string `json:"editor" bson:"editor"`
	EditorID primitive.ObjectID `json:"editor_id" bson:"editor_id"`
	Changes map[field]change `json:"changes" bson:"changes"`
}

func (self *logEvent) Log(ctx context.Context, collection string) *types.Deliverable {
	coll := db.Connection.Use(db.DefaultDatabase, collection)
	res, err := coll.InsertOne(ctx, self); if err != nil {
		//failed to log
		return self.failed()
	}

	uid, uidErr := UID.TryFromInterface(res.InsertedID); if uidErr != nil {
		//failed to log
		return self.failed()
	}

	var savedLogEvent loggedEvent
	err = coll.FindOne(ctx, bson.D{{"_id", uid.Oid()}}).Decode(&savedLogEvent)
	if err != nil {
		//failed to fetch log
		return self.failed()
	}
	return &savedLogEvent
}

func (self *logEvent) failed() *loggedEvent {
	return &loggedEvent{
		ID: primitive.ObjectID.NilObjectID,
		EventType: self.EventType,
		Timestamp: self.Timestamp,
		Editor: self.Editor,
		EditorID: self.EditorID,
		Changes: self.Changes,
	}
}
