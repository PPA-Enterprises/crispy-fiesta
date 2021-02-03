package eventlog

import (
	"PPA"
	"context"
	"fmt"
	"pkg/common/mongo"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	db *mongo.DBConnection
}

func New(db *mongo.DBConnection) Service {
	return Service{db}
}

func (s Service) structToMap(in interface{}, tag string) (PPA.EventMap, error) {
	out := make(map[PPA.Field]interface{})

	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// only take structs
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("structToMap only accepts structs; got %T", val)
	}

	typ := val.Type()
	for i := 0; i<val.NumField(); i++ {
		structField := typ.Field(i)
		if tagv := structField.Tag.Get(tag); tagv != "" {
			out[tagv] = val.Field(i).Interface()
		}
	}
	return out, nil
}

func (s Service) GenerateEvent(in interface{}, tag string) PPA.EventMap {
	ev, _ := s.structToMap(in, tag)
	return ev
}

func (s Service) LogCreated(ctx context.Context, data PPA.EventMap, editor PPA.Editor) PPA.LogEvent {
	changes := make(PPA.ChangesMap)

	for k, v := range data {
		changes[k] = PPA.Change{Old: nil, New: v}
	}
	fmt.Println(changes)

	event := PPA.LogEvent {
		ID: primitive.NewObjectID(),
		EventType: PPA.Created,
		Timestamp: time.Now(),
		Editor: editor.Name,
		EditorID: editor.OID,
		Changes: changes,
	}

	for s.oidExists(ctx, event.ID, editor.Collection) {
		event.ID = primitive.NewObjectID()
	}

	if !s.log(ctx, editor.Collection, &event) {
		return s.failed(&event)
	}
	return event
}

func (s Service) LogUpdated(ctx context.Context, prev PPA.EventMap, next PPA.EventMap, editor PPA.Editor) PPA.LogEvent {
	changes := make(PPA.ChangesMap)
	for k, v := range next {
		if !reflect.DeepEqual(next[k], prev[k]) {
			changes[k] = PPA.Change{Old: prev[k], New: v}
		}
	}

	event := PPA.LogEvent {
		ID: primitive.NewObjectID(),
		EventType: PPA.Edited,
		Timestamp: time.Now(),
		Editor: editor.Name,
		EditorID: editor.OID,
		Changes: changes,
	}

	for s.oidExists(ctx, event.ID, editor.Collection) {
		event.ID = primitive.NewObjectID()
	}

	if !s.log(ctx, editor.Collection, &event) {
		return s.failed(&event)
	}
	return event
}

func (s Service) LogDeleted(ctx context.Context, editor PPA.Editor) PPA.LogEvent {
	event := PPA.LogEvent {
		ID: primitive.NewObjectID(),
		EventType: PPA.Deleted,
		Timestamp: time.Now(),
		Editor: editor.Name,
		EditorID: editor.OID,
		Changes: PPA.ChangesMap{},
	}

	for s.oidExists(ctx, event.ID, editor.Collection) {
		event.ID = primitive.NewObjectID()
	}

	if !s.log(ctx, editor.Collection, &event) {
		return s.failed(&event)
	}
	return event
}

func (s Service) log(ctx context.Context, collection string, event *PPA.LogEvent) bool {
	coll := s.db.Use(collection)
	_, err := coll.InsertOne(ctx, event); if err != nil {
		return false
	}
	event.Persisted = true
	filter := bson.D{{"_id", event.ID}}
	updateDoc := bson.D{{"$set", event}}

	var oldDoc PPA.LogEvent
	err = coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	return err == nil
}

func (s Service) failed(event *PPA.LogEvent) PPA.LogEvent {
	return PPA.LogEvent{
		ID: primitive.NilObjectID,
		EventType: event.EventType,
		Timestamp: event.Timestamp,
		Editor: event.Editor,
		EditorID: event.EditorID,
		Persisted: false,
		Changes: event.Changes,
	}
}

func (s Service) oidExists(ctx context.Context, oid primitive.ObjectID, collection string) bool {
	coll := s.db.Use(collection)

	var inserted PPA.LogEvent
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}
