package jobs

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	errors "internal/common"
	"internal/db"
	"internal/uid"
)

type jobModel struct {
	ID              primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientName      string             `json:"client_info"bson:"client_info"`
	CarInfo         string             `json:"car_info"bson:"car_info"`
	AppointmentInfo string             `json:"appointment_info"bson:"appointment_info"`
	Notes           string             `json:"notes"bson:"notes"`
}

func fromSubmitJobCmd(data *submitJobCmd) *jobModel {
	return &jobModel{
		ClientName: data.ClientName,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo,
		Notes: data.Notes
}

func (self *job) create(ctx context.Context) (UID.ID, *errors.ResponseError) {
	//coll := db.Connection().Use(db.DefaultDatabase, "job")
	session, err := db.Connection().Session(); if err != nil {
		//failed to make a session
	}
	defer session.EndSession(ctx)

	session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		//if client exists, get it
		//if not, create a client
		//create job, update client array
	})

}
