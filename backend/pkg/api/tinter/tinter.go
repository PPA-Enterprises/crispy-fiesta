package tinter

import(
	"PPA"
	"context"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	Conflict = http.StatusConflict
	NotFound = http.StatusNotFound
	Collection = "tinters"
	EventTag = "m"
)

var OidNotFound = PPA.NewAppError(NotFound, "Does not exist")

type Update struct {
	Name string
	Phone string
}

func (t Tinter) Create(c *gin.Context, req PPA.Tinter, editor PPA.Editor) (*PPA.Tinter, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	// ensure OID is unique
	for t.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}

	created, err := t.tdb.Create(t.db, ctx, &req); if err != nil {
		return nil, err
	}

	created.AppendLog(t.eventLogger.LogCreated(ctx, t.eventLogger.GenerateEvent(created, EventTag), editor))
	t.tdb.LogEvent(t.db, ctx, created)
	return created, nil
}

func (t Tinter) ViewById(c *gin.Context, id string) (*PPA.Tinter, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, OidNotFound
	}

	return t.tdb.ViewById(t.db, ctx, oid)
}

func (t Tinter) List(c *gin.Context, opts PPA.BulkFetchOptions) (*[]PPA.Tinter, error) {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	return t.tdb.List(t.db, ctx, opts)
}

func (t Tinter) ViewByPhone(c *gin.Context, phone string) (*PPA.Tinter, error) {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	return t.tdb.ViewByPhone(t.db, ctx, phone)
}

func (t Tinter) Delete(c *gin.Context, id string, editor PPA.Editor) error {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return OidNotFound
	}

	tinter, err := t.tdb.ViewById(t.db, ctx, oid); if err != nil {
		return OidNotFound
	}
	// delete jobs
	//cl.deletejobs(ctx, client.Jobs, editor)

	tinter.AppendLog(t.eventLogger.LogDeleted(ctx, editor))
	t.tdb.LogEvent(t.db, ctx, tinter)
	return t.tdb.Delete(t.db, ctx, oid)

}

func (t Tinter) Update(c *gin.Context, req Update, id string, editor PPA.Editor) (*PPA.Tinter, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	// Note that len(nil) is 0
	if len(req.Phone) > 0 {
		fetched, _ := t.tdb.ViewByPhone(t.db, ctx, req.Phone)
		if fetched != nil {
			if fetched.ID.Hex() != id {
				return nil, PPA.NewAppError(Conflict, "Phone number already already in use")
			}
		}
	}

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, OidNotFound
	}

	oldDoc, err := t.tdb.ViewById(t.db, ctx, oid); if err != nil {
		return nil, err
	}

	if err := t.tdb.Update(t.db, ctx, oid, &PPA.Tinter {
		ID: primitive.NilObjectID,
		Name: req.Name,
		Phone: req.Phone,
	}); err != nil {
		return nil, err
	}

	updated, err := t.tdb.ViewById(t.db, ctx, oid); if err != nil {
		return nil, err
	}

	updated.AppendLog(t.eventLogger.LogUpdated(ctx,
		t.eventLogger.GenerateEvent(oldDoc, EventTag),
		t.eventLogger.GenerateEvent(updated, EventTag),
		editor))
	t.tdb.LogEvent(t.db, ctx, updated)

	// update client info on all corresponding jobs
	//cl.updateJobs(ctx, updated.Jobs, JobUpdate { Name: updated.Name, Phone: updated.Phone }, editor)
	return updated, nil
}
/*
func (t Tinter) PopulateJob(c *gin.Context, unpopClient *PPA.Client) (*PopulatedClient, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	jobs, err := cl.cdb.Populate(cl.db, ctx, unpopClient.Jobs); if err != nil {
		return nil, err
	}
	return &PopulatedClient {
		ID: unpopClient.ID,
		Name: unpopClient.Name,
		Phone: unpopClient.Phone,
		Jobs: jobs,
		History: unpopClient.History,
	}, nil
}

func (t Tinter) PopulateJobs(c *gin.Context, unpopClients *[]PPA.Client) (*[]PopulatedClient, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	var popClients = make([]PopulatedClient, 0, len(*unpopClients))

	for _, unpop := range *unpopClients {
		jobs, err := cl.cdb.Populate(cl.db, ctx, unpop.Jobs); if err != nil {
			//return nil, err
			//just skip it???
		}
		popClients = append(popClients, PopulatedClient {
			ID: unpop.ID,
			Name: unpop.Name,
			Phone: unpop.Phone,
			Jobs: jobs,
			History: unpop.History,
		})
	}
	return &popClients, nil
}*/

func (t Tinter) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := t.db.Use(Collection)

	var inserted PPA.Tinter
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}

/*
func (t Tinter) deletejobs(ctx context.Context, oids []primitive.ObjectID, editor PPA.Editor) {
	for _, oid := range oids {
		deleted, _ := t.jdb.ViewById(t.db, ctx, oid)
		if deleted != nil {
			deleted.AppendLog(t.eventLogger.LogDeleted(ctx, editor))
			t.jdb.LogEvent(t.db, ctx, deleted)
		}
		err := t.jdb.Delete(t.db, ctx, oid); if err != nil {
			// do nothing
		}
	}
}

// How strong of a consistency garuntee do we want????
func (t Tinter) updateJobs(ctx context.Context, oids []primitive.ObjectID, update JobUpdate, editor PPA.Editor) {
	if len(oids) <= 0 { return } // Note that len(nil) is 0

	for _, oid := range oids {
		oldDoc, _ := cl.jdb.ViewById(cl.db, ctx, oid)

		if err := cl.jdb.Update(cl.db, ctx, oid, &PPA.Job {
			ClientName: update.Name,
			ClientPhone: update.Phone,
		}); err != nil {
			// do nothing or fail, depends
		}

		newDoc, _ := cl.jdb.ViewById(cl.db, ctx, oid)

		if newDoc != nil && oldDoc != nil {
			newDoc.AppendLog(cl.eventLogger.LogUpdated(ctx,
				cl.eventLogger.GenerateEvent(oldDoc, EventTag),
				cl.eventLogger.GenerateEvent(newDoc, EventTag),
				editor))
			cl.jdb.LogEvent(cl.db, ctx, newDoc)
		}
	}
}*/
