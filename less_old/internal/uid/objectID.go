package UID
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"internal/common/errors"
)

type ID interface {
	String() string
	Oid() primitive.ObjectID
}

type uid struct {
	objectID primitive.ObjectID
}

func (self *uid) String() string {
	return self.objectID.Hex()
}

func (self *uid) Oid() primitive.ObjectID {
	return self.objectID
}

func TryFromInterface(id interface{}) (ID, *errors.ResponseError) {
	if res, ok := id.(primitive.ObjectID); ok {
		return &uid{objectID: res}, nil
	}
	return nil, errors.UidTypeAssertionError()
}

func FromOid(oid primitive.ObjectID) ID {
	return &uid{objectID: oid}
}
