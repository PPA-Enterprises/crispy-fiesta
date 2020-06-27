package models

import (
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID bson.ObjectId `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientInfo string `json:"client_info"bson:"client_info"`
	CarInfo string `json:"car_info"bson:"car_info"`
	AppointmentInfo string `json:"appointment_info"bson:"appointment_info"`
}

func FromSubmitJobCmd(data forms.SubmitJobCmd) *Job {
	j := Job{
		ClientInfo: data.ClientInfo,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo
	}
	return &j;
}


func (job *Job) PersistJob() primitive.ObjectId {
	coll := dbConnect.Use(databaseName, "job")
	res, err := coll.collection.InsertOne(context.Background(), job)
	if err != nil {
		panic(err)
	}

	return res.InsertedID
}
