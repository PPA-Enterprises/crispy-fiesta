package clients

import (
	"context"
	"internal/clients/types"
	"internal/common/errors"
	"internal/db"
	"internal/uid"
	"internal/event_log"
	eventLogTypes "internal/event_log/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type newClient struct {
	Name  string               `json:"name" bson:"name"`
	Phone string               `json:"phone" bson:"phone"`
	Jobs  []primitive.ObjectID `json:"jobs" bson:"jobs"`
}

func fromCreateClientCmd(data *createClientCmd) *newClient {
	return &newClient{
		Name: data.Name,
		Phone: data.Phone,
		Jobs: make([]primitive.ObjectID, 0),
	}
}

func NewClient(name, phone string) types.Client {
	return &clientModel{
		ID:    primitive.NewObjectID(),
		Name:  name,
		Phone: phone,
		Jobs:  []primitive.ObjectID{},
	}
}

func (self *newClient) createUniq(ctx context.Context, editor *eventLogTypes.Editor) (UID.ID, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	exists := ClientByPhone(ctx, self.Phone)
	if exists != nil {
		return nil, errors.DoesNotExist()
	}

	res, err := coll.InsertOne(ctx, self)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	oid, idErr := UID.TryFromInterface(res.InsertedID); if idErr != nil {
		return nil, idErr
	}

	createdClient, fetchErr := clientByID(ctx, oid.String()); if fetchErr != nil {
		return nil, fetchErr
	}

	loggedClient := event_log.LogCreated(ctx, createdClient.logable(), editor)
	_ = appendLog(ctx, createdClient, loggedClient)
	return oid, nil
}

