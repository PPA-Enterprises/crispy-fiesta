package PPA
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientName string `json:"client_info" bson:"client_name,omitempty" m:"Client Name"`
	ClientPhone string `json:"client_phone" bson:"client_phone,omitempty" m:"Client Phone Number"`
	CarInfo string `json:"car_info" bson:"car_info,omitempty" m:"Car Information"`
	AppointmentInfo string `json:"appointment_info" bson:"appointment_info,omitempty" m:"Appointment Information"`
	Notes string `json:"notes" bson:"notes,omitempty" m:"Notes"`
	History []LogEvent `json:"history" bson:"history,omitempty"`
}

func (j *Job) AppendLog(event LogEvent) {
	j.History = append(j.History, event)
}
