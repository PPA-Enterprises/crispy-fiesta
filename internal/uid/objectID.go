package UID
import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	errors "internal/common"
)

type ID interface {
	String() string
}

type uid struct {
	objectID primitive.ObjectID
}

func (self uid) String() string {
	return self.objectID.Hex()
}

func IdFromInterface(id interface{}) (ID, *errors.ResponseError) {
	if res, ok := id.(primitive.ObjectID); ok {
		return uid{objectID: res}, nil
	}
	return nil, errors.UidTypeAssertionError()
}
