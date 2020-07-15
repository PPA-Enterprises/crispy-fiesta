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
	ObjectID primitive.ObjectID
}

func (id uid) String() string {
	return id.ObjectID.String()
}

func IdFromInterface(id interface{}) (ID, error) {
	res, err := id.(primitive.ObjectID)
	log.Print(err)
	if !err {
		return nil, errors.New("Failed to get ID")
	}
	return uid{ObjectID: res}, nil
}
