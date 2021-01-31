package job
import (
	"PPA"
	"context"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Update struct {
	ClientName string
	ClientPhone string
	CarInfo string
	AppointmentInfo string
	Notes string
}

type ClientUpdate struct {
	Name string
	phone string
}

func (j Job) Create(c *gin.Context, req PPA.Job) (*PPA.Job, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	for j.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}

	created, err := j.jdb.Create(j.db, ctx, &req); if err != nil {
		return nil, err
	}

	if !j.clientExists(ctx, created.ClientPhone) {
		// create the client
		if _, err := j.createClient(ctx, &PPA.Client {
			ID: primitive.NilObjectID,
			Name: created.ClientName,
			Phone: created.ClientPhone,
			Jobs: []primitive.ObjectID{},
		}); err != nil {
			return nil, err
		}
	}

	// append id to client
	if err := j.attatchJobToClient(ctx, created.ClientPhone, created.ID); err != nil {
		return nil, err
	}

	return created, nil
}

func (j Job) ViewById( c *gin.Context, id string) (*PPA.Job, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, PPA.InternalError
	}
	return j.jdb.ViewById(j.db, ctx, oid)
}

func (j Job) List(c *gin.Context) (*[]PPA.Job, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	return j.jdb.List(j.db, ctx)
}

func (j Job) Delete(c *gin.Context, id string) error {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return PPA.InternalError
	}
	return j.jdb.Delete(j.db, ctx, oid)
}

func (j Job) Update(c *gin.Context, req Update, id string) (*PPA.Job, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, PPA.InternalError
	}
	// TODO: update name and number on client side too

	if err := j.jdb.Update(j.db, ctx, oid, &PPA.Job {
		ID: primitive.NilObjectID,
		ClientName: req.ClientName,
		ClientPhone: req.ClientPhone,
		CarInfo: req.CarInfo,
		AppointmentInfo: req.AppointmentInfo,
		Notes: req.Notes,
	}); err != nil {
		return nil, err
	}

	/* TODO: Needs discussion...update client name/phone?
	clientUpdate := ClientUpdate { Name: req.ClientName, Phone: re.ClientPhone }
	if err := updateClient(ctx, clientUpdate); err != nil {
		return nil, err
	}*/

	return j.jdb.ViewById(j.db, ctx, oid)
}

func (j Job) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := j.db.Use("jobs")

	var inserted PPA.Job
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}

func (j Job) clientExists(ctx context.Context, phone string) bool {
	_, err := j.cdb.ViewByPhone(j.db, ctx, phone)
	return err == nil
}

func (j Job) createClient(ctx context.Context, client *PPA.Client) (*PPA.Client, error) {
	client.ID = primitive.NewObjectID()
	for j.clientOidExists(ctx, client.ID) {
		client.ID = primitive.NewObjectID()
	}
	return j.cdb.Create(j.db, ctx, client)
}

func (j Job) clientOidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := j.db.Use("clients")

	var inserted PPA.Job
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}

func (j Job) attatchJobToClient(ctx context.Context, phone string, jobOid primitive.ObjectID) error {
	client, err := j.cdb.ViewByPhone(j.db, ctx, phone); if err != nil {
		return err
	}

	client.AttatchJob(jobOid)
	updateErr := j.cdb.Update(j.db, ctx, client.ID, client)
	return updateErr
}

