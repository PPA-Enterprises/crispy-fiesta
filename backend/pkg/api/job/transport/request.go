package transport
import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pkg/api/job"
)

type submitJobRequest struct {
	Title string`json:"title" binding:"required"`
	ClientName string `json:"client_name" binding:"required"`
	ClientPhone string `json:"client_phone" binding:"required"`
	CarInfo string `json:"car_info" binding:"required"`
	AppointmentInfo string `json:"appointment_info" binding:"required"`
	Notes string `json:"notes"`
	StartTime string `json:"start" binding:"required"`
	EndTime string `json:"end" binding:"required"`
	Tag string `json:"tag" binding:"required"`
	Color PPA.CalendarMeta `json:"color" binding:"required"`
}

type updateRequest struct {
	Title string`json:"title,omitempty:`
	ClientName string `json:"client_name,omitempty"`
	ClientPhone string `json:"client_phone,omitempty"`
	CarInfo string `json:"car_info,omitempty"`
	AppointmentInfo string `json:"appointment_info,omitempty"`
	Notes string `json:"notes,omitempty"`
	StartTime string `json:"start,omitempty"`
	EndTime string `json:"end,omitempty"`
	Tag string `json:"tag,omitempty"`
	Color *PPA.CalendarMeta `json:"color,omitempty"`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) job.Update {
	return job.Update {
		Title: data.Title,
		ClientName: data.ClientName,
		ClientPhone: data.ClientPhone,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo,
		Notes: data.Notes,
		StartTime: data.StartTime,
		EndTime: data.EndTime,
		Tag: data.Tag,
		Color: data.Color,
	}
}

func (h HTTP) fromSubmitJobRequest(data *submitJobRequest) PPA.Job {
	return PPA.Job {
		ID: primitive.NewObjectID(),
		Title: data.Title,
		ClientName: data.ClientName,
		ClientPhone: data.ClientPhone,
		CarInfo: data.CarInfo,
		AppointmentInfo: data.AppointmentInfo,
		Notes: data.Notes,
		StartTime: data.StartTime,
		EndTime: data.EndTime,
		Color: &data.Color,
		Tag: data.Tag,
	}
}
