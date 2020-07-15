package UID
import (
	"errors"
	"log"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func IdFromInterface(id interface{}) (ID, error) {
	res, err := id.(primitive.ObjectID)
	log.Print(err)
	if !err {
		return nil, errors.New("Failed to get ID")
	}
	return uid{objectID: res}, nil
}
