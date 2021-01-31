package transport
import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pkg/api/job"
)

type submitJobRequest struct {
	ClientName string `json:"client_name" binding:"required"`
	ClientPhone string `json:"client_phone" binding:"required"`
	CarInfo string `json:"car_info" binding:"required"`
	AppointmentInfo string `json:"appointment_info" binding:"required"`
	Notes string `json:"notes"`
}

type updateRequest struct {
	ClientName string `json:"client_name,omitempty"`
	ClientPhone string `json:"client_phone,omitempty"`
	CarInfo string `json:"car_info,omitempty"`
	AppointmentInfo string `json:"appointment_info,omitempty"`
	Notes string `json:"notes,omitempty"`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) job.Update {
	return job.Update {
		ClientName: data.ClientName,
		ClientPhone: data.ClientPhone,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo,
		Notes: data.Notes,
	}
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