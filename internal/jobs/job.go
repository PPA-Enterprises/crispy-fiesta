package jobs

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"internal/common/errors"
	"internal/db"
	"internal/uid"
	clientTypes "internal/clients/types"
	eventLogTypes "internal/event_log/types"
	"internal/event_log"
	"internal/clients"
	"internal/jobs/types"
)

type jobModel struct {
	ID              primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientName      string             `json:"client_info"bson:"client_info"`
	ClientPhone		string             `json:"client_phone"bson:"client_phone"`
	CarInfo         string             `json:"car_info"bson:"car_info"`
	AppointmentInfo string             `json:"appointment_info"bson:"appointment_info"`
	Notes           string             `json:"notes"bson:"notes"`
	Log				[]eventLogTypes.NormalizedLoggedEvent `json:"log"bson:"log"`
}

func fromSubmitJobCmd(data *submitJobCmd) *jobModel {
	return &jobModel{
		ID: primitive.NewObjectID(),
		ClientName: data.ClientName,
		ClientPhone: data.ClientPhone,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo,
		Notes: data.Notes,
	}
}

/*
// TODO: Mongo Transactions can only be done on a Mongo cluster such as a replica set
// Need additional infrastructure to support this

func (self *jobModel) create(ctx context.Context) (UID.ID, *errors.ResponseError) {
	session, err := db.Connection().Session(); if err != nil {
		//failed to make a session
		return nil, errors.DatabaseError(err)
	}
	defer session.EndSession(ctx)

	//transaction
	res, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		//if client exists, get it
		var client clients.Client
		client = clients.ClientByPhone(sessionCtx, self.ClientPhone)

		//if not, create a client
		if client == nil {
			client = clients.NewClient(self.ClientName, self.ClientPhone)
		}
		//update client array with job
		client.AttatchJobID(self.ID)

		//save job and client. Need Put for both to remain idempotent
		err := client.Put(sessionCtx); if err != nil {
			return nil, err
		}

		err = self.put(sessionCtx); if err != nil {
			return nil, err
		}

		return UID.FromOid(self.ID), nil
	})

	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	if id, ok := res.(UID.ID); ok {
		return id, nil
	}
	return nil, errors.UidTypeAssertionError()
}*/

func (self *jobModel) create(ctx context.Context, editor *eventLogTypes.Editor) (UID.ID, *errors.ResponseError) {
	//if client exists, get it
	var client clientTypes.Client
	client = clients.ClientByPhone(ctx, self.ClientPhone)

	//if not, create a client
	if client == nil {
		client = clients.NewClient(self.ClientName, self.ClientPhone)
	}

	//update client array with job
	err := client.AttatchJobID(ctx, self.ID, editor); if err != nil {
		return nil, err
	}

	err = self.put(ctx, true); if err != nil {
		return nil, err
	}

	loggedJob := event_log.LogCreated(ctx, self.logable(), editor)
	_ = appendLog(ctx, self, loggedJob)
	return UID.FromOid(self.ID), nil
}

func (self *jobModel) put(ctx context.Context, upsert bool) *errors.ResponseError {
	coll := db.Connection().Use(db.DefaultDatabase, "jobs")
	opts := options.FindOneAndReplace()
	opts = opts.SetUpsert(upsert)

	err := coll.FindOneAndReplace(ctx, bson.D{{"_id", self.ID}}, self, opts).Err()
	if err == mongo.ErrNoDocuments {
		if upsert {
			return nil
		} else {
			return errors.PutFailed(err)
		}
	}

	if err != nil {
		return errors.PutFailed(err)
	}
	return nil
}

func (self *jobModel) logable() *types.LogableJob {
	return &types.LogableJob {
		ID: self.ID.Hex(),
		ClientName: self.ClientName,
		ClientPhone: self.ClientPhone,
		CarInfo: self.CarInfo,
		AppointmentInfo: self.AppointmentInfo,
		Notes: self.Notes,
	}
}

func jobByID(ctx context.Context, id string) (*jobModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "jobs")

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, errors.InvalidOID()
	}

	var foundJob jobModel
	err = coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&foundJob)
	if err != nil {
		return nil, errors.DoesNotExist()
	}
	return &foundJob, nil
}

func deleteJobFromClient(ctx context.Context, jobID, clientID string) *errors.ResponseError {
	//TODO: delete job from client

	if err := DeleteByID(ctx, jobID); err != nil {
		return err
	}
	return nil
}

func DeleteByID(ctx context.Context, id string) *errors.ResponseError {
	coll := db.Connection().Use(db.DefaultDatabase, "deleted_jobs")
	job, err := jobByID(ctx, id); if err != nil {
		return err
	}
	_, insertErr := coll.InsertOne(ctx, job); if err != nil {
		return errors.DatabaseError(insertErr)
	}

	coll = db.Connection().Use(db.DefaultDatabase, "jobs")
	_, delErr := coll.DeleteOne(ctx, bson.D{{"_id", job.ID}})
	if delErr != nil {
		return errors.DatabaseError(delErr)
	}
	return nil
}
