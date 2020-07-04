package models

import (
	"context"

	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Job struct {
	ID              primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientName      string             `json:"client_info"bson:"client_info"`
	CarInfo         string             `json:"car_info"bson:"car_info"`
	AppointmentInfo string             `json:"appointment_info"bson:"appointment_info"`
	Notes           string             `json:"notes"bson:"notes"`
}

type JobModel struct{}

func (createJob *JobModel) CreateJob(data forms.SubmitJobCmd) ([2]*mongo.InsertOneResult, error) {
	var err error
	var client *mongo.Client
	var jobsCollection *mongo.Collection
	var clientCollection *mongo.Collection
	var ctx = context.Background()
	var result [2]*mongo.InsertOneResult
	var session mongo.Session
	client = dbConnect.Client
	// defer client.Disconnect(ctx)
	jobsCollection = client.Database("PPA").Collection("job")
	clientCollection = client.Database("PPA").Collection("client")

	if session, err = client.StartSession(); err != nil {
		panic(err)
	}
	if err = session.StartTransaction(); err != nil {
		panic(err)
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if result[0], err = jobsCollection.InsertOne(ctx, bson.D{
			{"carInfo", data.CarInfo},
			{"appointmentInfo", data.AppointmentInfo},
			{"notes", data.Notes},
		}); err != nil {
			panic(err)
		}

		if err = session.CommitTransaction(sc); err != nil {
			panic(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	var inProgress [1]primitive.ObjectID
	inProgress[0] = result[0].InsertedID.(primitive.ObjectID)
	var completed []primitive.ObjectID

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if result[1], err = clientCollection.InsertOne(ctx, bson.D{
			{"name", data.ClientName},
			{"inProgress", inProgress},
			{"completed", completed},
		}); err != nil {
			panic(err)
		}

		if err = session.CommitTransaction(sc); err != nil {
			panic(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	session.EndSession(ctx)

	return result, err
}

func (updateJob *JobModel) UpdateJob(data forms.UpdateJobCmd) (Job, error) {
	collection := dbConnect.Use(databaseName, "job")

	var returnedJob Job
	filter := bson.D{{"_id", data.ID}}
	// update := bson.D{{"$set", bson.D{
	// 	{"CarInfo", "newemail@example.com"}
	// }}}
	err := collection.FindOneAndUpdate(context.TODO(), filter, data).Decode(&returnedJob)

	return returnedJob, err

}
