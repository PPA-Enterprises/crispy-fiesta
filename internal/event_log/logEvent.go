package event_log

import (
	"time"
	"context"
	"internal/db"
	"internal/event_log/types"
	jobTypes "internal/jobs/types"
	"internal/uid"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type event = string
type field = string
const(
	created event = "Created"
	edited event = "Edited"
	deleted event = "Deleted"
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
		Persisted: false,
		Changes: self.Changes,
	}
}

func LogCreatedJob(ctx context.Context, job *jobTypes.LogableJob, editor *types.Editor) types.NormalizedLoggedEvent {
	changes, err := structToMap(job, "m"); if err != nil {
		return nil, err
	}
	return &types.loggedEvent {
		EventType: created,
		Timestamp: time.Now().Unix(),
		Editor: editor.Name,
		EditorID: editor.Oid,
		Changes: changes,
	}, nil
}
