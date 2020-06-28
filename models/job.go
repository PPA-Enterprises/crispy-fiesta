package models

import (
	"context"
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientInfo string `json:"client_info"bson:"client_info"`
	CarInfo string `json:"car_info"bson:"car_info"`
	AppointmentInfo string `json:"appointment_info"bson:"appointment_info"`
	Notes string `json:"notes"bson:"notes"`
}

func FromSubmitJobCmd(data forms.SubmitJobCmd) *Job {
	j := Job{
		ClientInfo: data.ClientInfo,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo}
	return &j;
}


func (job *Job) PersistJob() (*ID, error) {
	coll := dbConnect.Use(databaseName, "job")
	res, err := coll.InsertOne(context.Background(), job)
	if err != nil {
		return nil, err
	}

	return IdFromInterface(res.InsertedID)
}
