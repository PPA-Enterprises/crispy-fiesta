package clients

import (
	"context"
	"internal/clients/types"
	jobTypes "internal/jobs/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
func normalize(j []jobTypes.Job) []primitive.ObjectID {
	oids := make([]primitive.ObjectID, 0)
	for _, job := range j {
		oids = append(oids, job.ID)
	}
	return oids
}

func populateClients(ctx context.Context, clients []clientModel) []types.PopulatedClientModel {
	populatedClients := make([]types.PopulatedClientModel, 0, len(clients))
	for _, c := range clients {
		client, err := c.Populate(ctx)
		//just skip bad records
		if err == nil {
			populatedClients = append(populatedClients, *client)
		}
	}
	return populatedClients
}
