package transport
import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type submitJobRequest struct {
	ClientName string `json:"client_name" binding:"required"`
	ClientPhone string `json:"client_phone" binding:"required"`
	CarInfo string `json:"car_info" binding:"required"`
	AppointmentInfo string `json:"appointment_info" binding:"required"`
	Notes string `json:"notes"`
}

func (h HTTP) fromSubmitJobRequest(data *submitJobRequest) PPA.Job {
	return PPA.Job {
		ID: primitive.NewObjectID(),
		ClientName: data.ClientName,
		ClientPhone: data.ClientPhone,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo,
		Notes: data.Notes,
	}
}
