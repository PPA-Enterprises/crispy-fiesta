package transport

import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pkg/api/client"
)

type createClientRequest struct {
	Name string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

func (h HTTP) fromCreateClientRequest(data *createClientRequest) PPA.Client {
	return PPA.Client {
		Name: data.Name,
		Phone: data.Phone,
	}
}