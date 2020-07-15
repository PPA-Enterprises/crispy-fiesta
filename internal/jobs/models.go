package jobs

import (
	//"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
)

type job struct {
	ID              primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientName      string             `json:"client_info"bson:"client_info"`
	CarInfo         string             `json:"car_info"bson:"car_info"`
	AppointmentInfo string             `json:"appointment_info"bson:"appointment_info"`
	Notes           string             `json:"notes"bson:"notes"`
}

