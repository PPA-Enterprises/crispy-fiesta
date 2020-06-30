package models

import (
	"context"

	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	ID              primitive.ObjectID `json:"_id,omitempty"bson:"_id,omitempty"`
	ClientInfo      string             `json:"client_info"bson:"client_info"`
	CarInfo         string             `json:"car_info"bson:"car_info"`
	AppointmentInfo string             `json:"appointment_info"bson:"appointment_info"`
	Notes           string             `json:"notes"bson:"notes"`
}

type JobModel struct{}

func FromSubmitJobCmd(data forms.SubmitJobCmd) *Job {
	j := Job{
		ClientInfo:      data.ClientInfo,
		CarInfo:         data.CarInfo,
		AppointmentInfo: data.AppointmentInfo}
	return &j
}

func (job *Job) PersistJob() (*ID, error) {
	/*coll := dbConnect.Use(databaseName, "job")
	res, err := coll.InsertOne(context.Background(), job)
	if err != nil {
		return nil, err
	}*/

	session, err := dbConnect.Session()
	if err != nil {
		//500
	}
	defer session.EndSession(context.TODO())
	sessCtx := mongo.NewSessionContext(context.TODO(), session)

	if err = session.StartTransaction(); err != nil {
		panic(err)
		//return err
	}
	coll := dbConnect.Use(databaseName, "job")
	res, err := coll.InsertOne(sessCtx, job)
	if err != nil {
		if transErr = session.AbortTransaction(); transErr != nil {
			//no transaction to abort
			//500
		}
		return nil, err
	}

	id, err := IdFromInterface(res.InsertedID)
	if err != nil {
		if transErr = session.AbortTransaction(); transErr != nil {
			//no transaction to abort
			//500
		}
		return nil, err
	}

	// see if client exists
	//yes, add id
	//no, create it and add id

	return id, nil
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
