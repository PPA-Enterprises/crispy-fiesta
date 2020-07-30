package types

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	errors "internal/common"
)

type Client interface {
	AttatchJobID(primitive.ObjectID)
	Put(ctx context.Context) *errors.ResponseError
}
