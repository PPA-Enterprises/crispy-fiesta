package clients

import (
	"internal/clients/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type joblessClient struct {
	ID primitive.ObjectID	`json:"_id"`
	Name string				`json:"name"`
	Phone string			`json:"phone"`
	Jobs []string			`json:"jobs"`
}

func emptyJobsClient(c *types.PopulatedClientModel) *joblessClient {
	return &joblessClient {
		ID: c.ID,
		Name: c.Name,
		Phone: c.Phone,
		Jobs: make([]string, 0),
	}
}
