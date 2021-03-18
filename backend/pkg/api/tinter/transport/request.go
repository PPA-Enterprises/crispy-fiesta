package transport

import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pkg/api/tinter"
)

type createTinterRequest struct {
	Name string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type updateRequest struct {
	Name string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) tinter.Update {
	return client.Update {
		Name: data.Name,
		Phone: data.Phone,
	}
}

func (h HTTP) fromCreateClientRequest(data *createtinterRequest) PPA.Tinter {
	return PPA.Tinter {
		ID: primitive.NewObjectID(),
		Name: data.Name,
		Phone: data.Phone,
		//Jobs: []primitive.ObjectID{},
	}
}
