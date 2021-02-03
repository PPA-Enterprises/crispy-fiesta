package PPA
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientName string `json:"client_info" bson:"client_name,omitempty"`
	ClientPhone string `json:"client_phone" bson:"client_phone,omitempty"`
	CarInfo string `json:"car_info" bson:"car_info,omitempty"`
	AppointmentInfo string `json:"appointment_info" bson:"appointment_info,omitempty"`
	Notes string `json:"notes" bson:"notes,omitempty"`
	History []LogEvent `json:"history" bson:"history,omitempty"`
}

func (j *Job) AppendLog(event LogEvent) {
	j.History = append(j.History, event)
}
