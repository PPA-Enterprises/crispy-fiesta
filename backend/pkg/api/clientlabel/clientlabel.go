package clientlabel

import (
	"context"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"PPA"
)

const (
	NotFound = http.StatusNotFound
)

var OidNotFound = PPA.NewAppError(NotFound, "Does not exist")
const (
	Collection = "clientlabels"
	EventTag = "m"
)

type Update struct {
	LabelName string
}

func (cl ClientLabel) Create(c *gin.Context, req PPA.ClientLabel, editor PPA.Editor) (*PPA.ClientLabel, error) {
	//additional security stuff like if user is allowed to do this
	/*if err := self.rbac.AccountCreate(c); err != nil {
		return PPA.User{}, err
	}*/
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	// ensure OID is unique
	for cl.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}

	created, err := cl.cldb.Create(cl.db, ctx, &req); if err != nil {
		return nil, err
	}

	created.AppendLog(cl.eventLogger.LogCreated(ctx, cl.eventLogger.GenerateEvent(created, EventTag), editor))
	cl.cldb.LogEvent(cl.db, ctx, created)

	return created, nil
}

func (cl ClientLabel) ViewById(c *gin.Context, id string) (*PPA.ClientLabel, error) {
	//additional security stuff
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, OidNotFound
	}

	return cl.cldb.ViewById(cl.db, ctx, oid)
}

func (cl ClientLabel) List(c *gin.Context) (*[]PPA.ClientLabel, error) {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	//Populate anything?
	return cl.cldb.List(cl.db, ctx)
}

func (cl ClientLabel) Delete(c *gin.Context, id string, editor PPA.Editor) error {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return OidNotFound
	}

	oldDoc, err := cl.cldb.ViewById(u.db, ctx, oid); if err != nil {
		return OidNotFound
	}

	oldDoc.AppendLog(cl.eventLogger.LogDeleted(ctx, editor))
	cl.cldb.LogEvent(cl.db, ctx, oldDoc)

	return cl.cldb.Delete(cl.db, ctx, oid)
}


func (cl ClientLabel) Update(c *gin.Context, req Update, id string, editor PPA.Editor) (*PPA.ClientLabel, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, OidNotFound
	}

	oldDoc, err := cl.cldb.ViewById(cl.db, ctx, oid); if err != nil {
		return nil, OidNotFound
	}

	if err := cl.cldb.Update(cl.db, ctx, oid, &PPA.ClientLabel {
		ID: primitive.NilObjectID,
		LabelName: req.LabelName,
	}); err != nil {
		return nil, err
	}

	updated, err := cl.cldb.ViewById(cl.db, ctx, oid); if err != nil {
		return nil, PPA.InternalError
	}

	updated.AppendLog(cl.eventLogger.LogUpdated(ctx,
		cl.eventLogger.GenerateEvent(oldDoc, EventTag),
		cl.eventLogger.GenerateEvent(updated, EventTag),
		editor))
	cl.cldb.LogEvent(cl.db, ctx, updated)

	return updated, nil
}


func (cl ClientLabel) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := cl.db.Use(Collection)

	var inserted PPA.ClientLabel
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}
