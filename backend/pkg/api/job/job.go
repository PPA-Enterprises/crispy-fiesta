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
	ID primitive.ObjectID
	ClientName string
	ClientPhone string
	CarInfo string
	AppointmentInfo string
	Notes string
}

func (j Job) Create(c *gin.Context, req PPA.Job) (*PPA.Job, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	for j.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}
	return j.jdb.Create(j.db, ctx, &req)
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
		return PPA.InternalError
	}

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
	return j.jdb.ViewById(j.db, ctx, oid)
}

func (j Job) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := j.db.Use("job")

	var inserted PPA.Job
	err := Coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}
