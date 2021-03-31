package transport

import (
	"PPA"
	"pkg/api/job"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type submitJobRequest struct {
	Title           string				`json:"title" binding:"required"`
	ClientName      string				`json:"client_name" binding:"required"`
	ClientPhone     string				`json:"client_phone" binding:"required"`
	AssignedWorker	string				`json:"assigned_worker"`
	CarInfo         string				`json:"car_info" binding:"required"`
	Notes           string				`json:"notes"`
	StartTime       string				`json:"start" binding:"required"`
	EndTime         string				`json:"end" binding:"required"`
	Tag             string				`json:"tag" binding:"required"`
	Color           PPA.CalendarMeta	`json:"color" binding:"required"`
}

type updateRequest struct {
	Title           string				`json:"title,omitempty:`
	ClientName      string				`json:"client_name,omitempty"`
	ClientPhone     string				`json:"client_phone,omitempty"`
	AssignedWorker	string				`json:"assigned_worker,omitempty"`
	CarInfo         string				`json:"car_info,omitempty"`
	Notes           string				`json:"notes,omitempty"`
	StartTime       string				`json:"start,omitempty"`
	EndTime         string				`json:"end,omitempty"`
	Tag             string				`json:"tag,omitempty"`
	Color           *PPA.CalendarMeta	`json:"color,omitempty"`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) job.Update {
	assignedWorkerOID, err := primitive.ObjectIDFromHex(data.AssignedWorker); if err != nil {
		assignedWorkerOID = primitive.NilObjectID
	}
	return job.Update{
		Title:           data.Title,
		ClientName:      data.ClientName,
		ClientPhone:     data.ClientPhone,
		AssignedWorker:  assignedWorkerOID,
		CarInfo:         data.CarInfo,
		Notes:           data.Notes,
		StartTime:       data.StartTime,
		EndTime:         data.EndTime,
		Tag:             data.Tag,
		Color:           data.Color,
	}
}

func (h HTTP) fromSubmitJobRequest(data *submitJobRequest) PPA.Job {
	assignedWorkerOID, err := primitive.ObjectIDFromHex(data.AssignedWorker); if err != nil {
		assignedWorkerOID = primitive.NilObjectID
	}
	return PPA.Job{
		ID:              primitive.NewObjectID(),
		Title:           data.Title,
		ClientName:      data.ClientName,
		ClientPhone:     data.ClientPhone,
		AssignedWorker:  assignedWorkerOID,
		CarInfo:         data.CarInfo,
		Notes:           data.Notes,
		StartTime:       data.StartTime,
		EndTime:         data.EndTime,
		Color:           &data.Color,
		Tag:             data.Tag,
	}
}
