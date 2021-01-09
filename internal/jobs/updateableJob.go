package jobs

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"internal/common/errors"
	"internal/db"
	"internal/event_log"
)

type updateableJob struct {
	ID              primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientName      string             `json:"client_info"bson:"client_info,omitempty"`
	ClientPhone		string             `json:"client_phone"bson:"client_phone,omitempty"`
	CarInfo         string             `json:"car_info"bson:"car_info,omitempty"`
	AppointmentInfo string             `json:"appointment_info"bson:"appointment_info,omitempty"`
	Notes           string             `json:"notes"bson:"notes,omitempty"`
}

func tryFromUpdateJobCmd(data *updateJobCmd, id string) (*updateableJob, *errors.ResponseError) {
	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, errors.InvalidOID()
	}

	return &updateableJob{
		ID: oid,
		ClientName: data.ClientName,
		ClientPhone: data.ClientPhone,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo,
		Notes: data.Notes,
	}, nil
}

func (self *updateableJob) Patch(ctx context.Context, upsert bool) (*jobModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "jobs")
	opts := options.FindOneAndUpdate().SetUpsert(upsert)

	filter := bson.D{{"_id", self.ID}}
	update := bson.D{{"$set", self}}
	var oldDocument jobModel
	var updatedDocument jobModel

	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&oldDocument)
	if err != nil {
		return nil, errors.PutFailed(err)
	}

	err = coll.FindOne(ctx, filter).Decode(&updatedDocument)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	loggedJob := event_log.LogUpdated(ctx, oldDocument.logable(), updatedDocument.logable(), editor)
	_ = appendLog(ctx, &updatedDocument, loggedJob)
	return &updatedDocument, nil
}
