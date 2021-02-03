package jobs
import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"internal/common/errors"
	"internal/db"
	eventLogTypes "internal/event_log/types"
)

func appendLog(ctx context.Context, job *jobModel, event *eventLogTypes.NormalizedLoggedEvent) *errors.ResponseError {
	//append log to job history
	job.Log = append(job.Log, *event)
	coll := db.Connection().Use(db.DefaultDatabase, "jobs")
	opts := options.FindOneAndUpdate().SetUpsert(false)

	filter := bson.D{{"_id", job.ID}}
	updateLog := bson.D{{"$set", bson.D{{"log", job.Log}}}}
	var updatedDocument jobModel

	err := coll.FindOneAndUpdate(ctx, filter, updateLog, opts).Decode(&updatedDocument)
	if err != nil {
		return errors.PutFailed(err)
	}

	err = coll.FindOne(ctx, filter).Decode(&updatedDocument)
	if err != nil {
		return errors.DatabaseError(err)
	}
	job.Log = updatedDocument.Log
	return nil
}
