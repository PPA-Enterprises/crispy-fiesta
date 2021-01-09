package event_log

import (
	"context"
	"time"
	"internal/db"
	"internal/uid"
	"internal/event_log/types"
	jobTypes "internal/jobs/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

type event = string
type field = string
const(
	created event = "Created"
	edited event = "Edited"
	deleted event = "Deleted"
)

type logEvent struct {
	EventType event `json:"event_type" bson:"event_type"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	Editor string `json:"editor" bson:"editor"`
	EditorID primitive.ObjectID `json:"editor_id" bson:"editor_id"`
	Changes map[field]types.Change `json:"changes" bson:"changes"`
}

func (self *logEvent) log(ctx context.Context, collection string) *loggedEvent {
	coll := db.Connection().Use(db.DefaultDatabase, collection)
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
	savedLogEvent.Persisted = true
	return &savedLogEvent
}

func (self *logEvent) failed() *loggedEvent {
	return &loggedEvent{
		ID: primitive.NilObjectID,
		EventType: self.EventType,
		Timestamp: self.Timestamp,
		Editor: self.Editor,
		EditorID: self.EditorID,
		Persisted: false,
		Changes: self.Changes,
	}
}

func LogCreated(ctx context.Context, data interface{}, editor *types.Editor) *types.NormalizedLoggedEvent {
	var changesMap map[string]interface{}
	var err error

	v, ok := data.(*jobTypes.LogableJob); if ok {
		changesMap, err = structToMap(v, "m"); if err != nil {
			return nil
		}
	}

	changes := make(map[field]types.Change)
	for key, value := range changesMap {
		changes[key] = types.Change{Old:nil, New:value}
	}

	event := &logEvent {
		EventType: created,
		Timestamp: time.Now().Unix(),
		Editor: editor.Name,
		EditorID: editor.Oid,
		Changes: changes,
	}
	return event.log(ctx, editor.Collection).normalize()
}
