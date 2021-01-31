package client

import(
	"PPA"
	"context"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

type Update struct {
	Name string
	Phone string
}

func (cl Client) Create(c *gin.Context, req PPA.Client) (*PPA.Client, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	// ensure OID is unique
	for cl.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}
	return cl.cdb.Create(cl.db, ctx, &req)
}

func (cl Client) ViewById(c *gin.Context, id string) (*PPA.Client, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, PPA.InternalError
	}

	return cl.cdb.ViewById(cl.db, ctx, oid)
}

func (cl Client) List(c *gin.Context) (*[]PPA.Client, error) {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	//Populate anything?
	return cl.cdb.List(cl.db, ctx)
}

func (cl Client) ViewByPhone(c *gin.Context, phone string) (*PPA.Client, error) {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	return cl.cdb.ViewByPhone(cl.db, ctx, phone)
}

func (cl Client) Delete(c *gin.Context, id string) error {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		//TODO: return an error about the object ID.
		//do this for anytime that I try to convert hex to OID
		return PPA.InternalError
	}
	// TODO: Delete all jobs assigned to client

	return cl.cdb.Delete(cl.db, ctx, oid)

}

func (cl Client) Update(c *gin.Context, req Update, id string) (*PPA.Client, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, PPA.InternalError
	}
	//TODO: Ensure updated phone number is unique

	if err := cl.cdb.Update(cl.db, ctx, oid, &PPA.Client {
		ID: primitive.NilObjectID,
		Name: req.Name,
		Phone: req.Phone,
	}); err != nil {
		return nil, err
	}
	return cl.cdb.ViewById(cl.db, ctx, oid)
}

func (cl Client) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := cl.db.Use("clients")

	var inserted PPA.Client
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}
