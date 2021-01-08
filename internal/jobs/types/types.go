package types

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//type Job interface {}
type Job struct {
	ID              primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientName      string             `json:"client_info"bson:"client_info"`
	ClientPhone		string             `json:"client_phone"bson:"client_phone"`
	CarInfo         string             `json:"car_info"bson:"car_info"`
	AppointmentInfo string             `json:"appointment_info"bson:"appointment_info"`
	Notes           string             `json:"notes"bson:"notes"`
}

type LogableJob struct {
	ID string
	ClientName string
	ClientPhone string
	CarInfo string
	AppointmentInfo string
	Notes string
}
