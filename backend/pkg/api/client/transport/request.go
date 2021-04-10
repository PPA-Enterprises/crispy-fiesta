package transport

import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	"pkg/api/client"
)

type createClientRequest struct {
	Name string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Labels []string `json:"labels"`
}

type updateRequest struct {
	Name string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type putLabelsRequest struct {
	Labels []string `json:"labels" binding:"required"`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) client.Update {
	return client.Update {
		Name: data.Name,
		Phone: data.Phone,
	}
}

func (h HTTP) tryFromCreateClientRequest(c *gin.Context, data *createClientRequest) (PPA.Client, error) {
	if len(data.Labels) <= 0 {
		return PPA.Client {
			ID: primitive.NewObjectID(),
			Name: data.Name,
			Phone: data.Phone,
			Labels: []primitive.ObjectID{},
			Jobs: []primitive.ObjectID{},
		}, nil
	}

	labelOIDs, err := h.service.FetchLabelOIDs(c, data.Labels); if err != nil {
		return PPA.Client{}, err
	}

	return PPA.Client {
		ID: primitive.NewObjectID(),
		Name: data.Name,
		Phone: data.Phone,
		Labels: labelOIDs,
		Jobs: []primitive.ObjectID{},
	}, nil
}
