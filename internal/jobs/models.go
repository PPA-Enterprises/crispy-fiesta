package jobs

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	errors "internal/common"
	"internal/db"
	"internal/uid"
	"internal/clients"
)

type Job interface {}

type jobModel struct {
	ID              primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientName      string             `json:"client_info"bson:"client_info"`
	ClientPhone		string             `json:"client_phone"bson:"client_phone"`
	CarInfo         string             `json:"car_info"bson:"car_info"`
	AppointmentInfo string             `json:"appointment_info"bson:"appointment_info"`
	Notes           string             `json:"notes"bson:"notes"`
}

func fromSubmitJobCmd(data *submitJobCmd) *jobModel {
	return &jobModel{
		ID: bson.NewObjectId()
		ClientName: data.ClientName,
		ClientPhone: data.ClientPhone,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo,
		Notes: data.Notes
}

func (self *job) create(ctx context.Context) (*UID.ID, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "jobs")
	session, err := db.Connection().Session(); if err != nil {
		//failed to make a session
		return nil, errors.DatabaseError(err)
	}
	defer session.EndSession(ctx)

	//ACID transaction
	res, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		//if client exists, get it
		var client *clients.Client
		client = clients.ClientByPhone(sessionCtx, self.ClientPhone)

		//if not, create a client
		if client == nil {
			client = clients.NewClient(self.ClientName, self.ClientPhone)
		}
		//update client array with job
		client.AttatchJobID(self.ID)

		//save job and client
		err := client.put(sessionCtx); if err != nil {
			return nil, err
		}

		err := job.put(sessionCtx); if err != nil {
			return nil, err
		}

		return UID.FromOid(self.ID), nil
	})

	if err != nil {
		nil, errors.DatabaseError(err)
	}
	return res, nil

}
func (self *job) put(ctx context.Context) errors.ResponseError {
	coll := db.Connection().Use(db.DefaultDatabase, "jobs")
	opts := options.FindOneAndReplace()
	opts = opts.SetUpsert(true)

	err := coll.FindOneAndReplace(ctx, bson.D{{"_id", self.ID}}).Err()
	if err == mongo.ErrNoDocuments {
		return nil
	}

	if err != nil {
		return errors.PutFailed(err)
	}
	return nil
}
