package transport

import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pkg/api/tinter"
)

type createTinterRequest struct {
	Name string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Jobs []string `json:"jobs"`
}

type updateRequest struct {
	Name string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) tinter.Update {
	return tinter.Update {
		Name: data.Name,
		Phone: data.Phone,
	}
}

func (h HTTP) tryFromCreateTinterRequest(data *createTinterRequest) (PPA.Tinter, error) {
	oids := make([]primitive.ObjectID, 0)
	if len(data.Jobs) > 0 {
		for _, id := range data.Jobs {
			oid, err := primitive.ObjectIDFromHex(id); if err != nil { return PPA.Tinter{}, err }
			oids = append(oids, oid)
		}
	}

	return PPA.Tinter {
		ID: primitive.NewObjectID(),
		Name: data.Name,
		Phone: data.Phone,
		Jobs: oids,
	}, nil
}
