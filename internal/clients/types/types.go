package types

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"internal/common/errors"
	jobTypes "internal/jobs/types"
)

type PopulatedClientModel struct {
	ID primitive.ObjectID `json:"_id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Jobs []jobTypes.Job `json:"jobs"`
}

type UnpopulatedClientModel struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	Jobs []primitive.ObjectID `json:"jobs" bson"jobs"`
}

type Client interface {
	AttatchJobID(primitive.ObjectID)
	Put(ctx context.Context, upsert bool) (*errors.ResponseError)
	Populate(ctx context.Context) (*PopulatedClientModel, *errors.ResponseError)
}

type DeliverableClient struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	Jobs []jobTypes.Job `json:"jobs" bosn:"jobs"`
}
