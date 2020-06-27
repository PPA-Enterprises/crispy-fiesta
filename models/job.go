package models

import (
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	ID bson.ObjectId `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientInfo string `json:"client_info"bson:"client_info"`
	CarInfo string `json:"car_info"bson:"car_info"`
	AppointmentInfo string `json:"appointment_info"bson:"appointment_info"`
}

type JobModel struct {}


