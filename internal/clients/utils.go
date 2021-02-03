package clients

import (
	"context"
	"internal/clients/types"
	"internal/db"
	"internal/common/errors"
	jobTypes "internal/jobs/types"
	eventLogTypes "internal/event_log/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
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

func appendLog(ctx context.Context, client *clientModel, event *eventLogTypes.NormalizedLoggedEvent) *errors.ResponseError {
	client.Log = append(client.Log, *event)
	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	opts := options.FindOneAndUpdate().SetUpsert(false)

	filter := bson.D{{"_id", client.ID}}
	updateLog := bson.D{{"$set", bson.D{{"log", client.Log}}}}
	var updatedDoc clientModel

	err := coll.FindOneAndUpdate(ctx, filter, updateLog, opts).Decode(&updatedDoc)
	if err != nil {
		return errors.PutFailed(err)
	}

	err = coll.FindOne(ctx, filter).Decode(&updatedDoc)
	if err != nil {
		return errors.DatabaseError(err)
	}
	client.Log = updatedDoc.Log
	return nil
}
