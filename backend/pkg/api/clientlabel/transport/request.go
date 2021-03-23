package transport

import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pkg/api/clientlabel"
)

type createRequest struct {
	LabelName string `json:"label_name" binding:"required,alphanum"`
}

type updateRequest struct {
	LabelName string `json:"label_name,omitempty"`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) clientlabel.Update {
	return clientlabel.Update {
		LabelName: data.LabelName,
	}
}

func (h HTTP) fromCreateRequest(data *createRequest) PPA.ClientLabel {
	oid := primitive.NewObjectID()

	return PPA.ClientLabel {
		ID: oid,
		LabelName: data.LabelName,
		Clients: make([]primitive.ObjectID, 0),
		History: make([]PPA.LogEvent, 0),
	}
}
