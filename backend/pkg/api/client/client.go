package client

import(
	"PPA"
	"fmt"
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
)

var OidNotFound = PPA.NewAppError(NotFound, "Does not exist")

type Update struct {
	Name string
	Phone string
}

type JobUpdate struct {
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
		return nil, OidNotFound
	}

	return cl.cdb.ViewById(cl.db, ctx, oid)
}

func (cl Client) List(c *gin.Context, opts PPA.BulkFetchOptions) (*[]PPA.Client, error) {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	return cl.cdb.List(cl.db, ctx, opts)
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
		return OidNotFound
	}

	client, err := cl.cdb.ViewById(cl.db, ctx, oid); if err != nil {
		return OidNotFound
	}
	// delete jobs
	cl.deletejobs(ctx, client.Jobs)

	return cl.cdb.Delete(cl.db, ctx, oid)

}

func (cl Client) Update(c *gin.Context, req Update, id string) (*PPA.Client, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	// Note that len(nil) is 0
	if len(req.Phone) > 0 {
		fetched, _ := cl.cdb.ViewByPhone(cl.db, ctx, req.Phone)
		fmt.Println(fetched)
		if fetched != nil {
			if fetched.ID.Hex() != id {
				return nil, PPA.NewAppError(Conflict, "Phone number already already in use")
			}
		}
	}

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, OidNotFound
	}

	if err := cl.cdb.Update(cl.db, ctx, oid, &PPA.Client {
		ID: primitive.NilObjectID,
		Name: req.Name,
		Phone: req.Phone,
	}); err != nil {
		return nil, err
	}

	updated, err := cl.cdb.ViewById(cl.db, ctx, oid); if err != nil {
		return nil, err
	}

	// update client info on all corresponding jobs
	cl.updateJobs(ctx, updated.Jobs, JobUpdate { Name: updated.Name, Phone: updated.Phone })
	return updated, nil
}

func (cl Client) PopulateJob(c *gin.Context, unpopClient *PPA.Client) (*PopulatedClient, error) {
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
	}, nil
}

func (cl Client) PopulateJobs(c *gin.Context, unpopClients *[]PPA.Client) (*[]PopulatedClient, error) {
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
		})
	}
	return &popClients, nil
}

func (cl Client) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := cl.db.Use("clients")

	var inserted PPA.Client
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}

func (cl Client) deletejobs(ctx context.Context, oids []primitive.ObjectID) {
	for _, oid := range oids {
		err := cl.jdb.Delete(cl.db, ctx, oid); if err != nil {
			// do nothing
		}
	}
}

// How strong of a consistency garuntee do we want????
func (cl Client) updateJobs(ctx context.Context, oids []primitive.ObjectID, update JobUpdate) {
	if len(oids) <= 0 { return } // Note that len(nil) is 0

	for _, oid := range oids {
		if err := cl.jdb.Update(cl.db, ctx, oid, &PPA.Job {
			ClientName: update.Name,
			ClientPhone: update.Phone,
		}); err != nil {
			// do nothing or fail, depends
		}
	}
}
