package PPA
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title,omitempty" m:"Title of Job"`
	StartTime string `json:"start" bson:"start_time,omitempty" m:"Start Time"`
	EndTime string `json:"end" bson:"end_time,omitempty" m:"End Time"`
	Tag string `json:"tag" bson:"tag,omitempty" m:"tag"`
	AssignedWorker primitive.ObjectID `json:"assigned_worker" bson:"assigned_worker,omitempty"`
	ClientName string `json:"client_name" bson:"client_name,omitempty" m:"Client Name"`
	ClientPhone string `json:"client_phone" bson:"client_phone,omitempty" m:"Client Phone Number"`
	CarInfo string `json:"car_info" bson:"car_info,omitempty" m:"Car Information"`
	Notes string `json:"notes" bson:"notes,omitempty" m:"Notes"`
	Color *CalendarMeta `json:"color" bson:"color,omitempty" m:"Calendar Colors"`
	History []LogEvent `json:"history" bson:"history,omitempty"`
}

type CalendarMeta struct {
	PrimaryColor string `json:"primary" bson:"primary_color,omitempty" m:"Calendar Primary Color"`
	SecondaryColor string `json:"secondary" bson:"secondary_color,omitempty" m:"Calendar Secondary Color"`
}

func (j *Job) AppendLog(event LogEvent) {
	j.History = append(j.History, event)
}
