package PPA
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StartTime string `json:"start" bson":start_time" m:"Start Time"`
	EndTime string `json:"end" bson:"end_time" m:"End Time"`
	Tag string `json:"tag" bson:"tag" m:"tag"`
	ClientName string `json:"client_info" bson:"client_name,omitempty" m:"Client Name"`
	ClientPhone string `json:"client_phone" bson:"client_phone,omitempty" m:"Client Phone Number"`
	CarInfo string `json:"car_info" bson:"car_info,omitempty" m:"Car Information"`
	AppointmentInfo string `json:"appointment_info" bson:"appointment_info,omitempty" m:"Appointment Information"`
	Notes string `json:"notes" bson:"notes,omitempty" m:"Notes"`
	Color CalendarMeta `json:"color" bson:"color" m:"Calendar Colors"`
	History []LogEvent `json:"history" bson:"history,omitempty"`
}

type CalendarMeta struct {
	PrimaryColor string `json:"primary" bson:"primary_color" m:"Calendar Primary Color"`
	SecondaryColor string `json:"secondary" bson:"secondary_color" m:"Calendar Secondary Color"`
}

func (j *Job) AppendLog(event LogEvent) {
	j.History = append(j.History, event)
}
