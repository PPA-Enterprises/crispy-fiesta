package client

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
	Collection = "clients"
	EventTag = "m"
)

var OidNotFound = PPA.NewAppError(NotFound, "Does not exist")

type PopulatedClient struct {
	ID primitive.ObjectID `json:"_id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Jobs []PPA.Job `json:"jobs"`
	Labels []string `json:"labels"`
	History []PPA.LogEvent `json:"history"`
}
type Update struct {
	Name string
	Phone string
}

type JobUpdate struct {
	Name string
	Phone string
}

// This is used for Logging only
type logableLabeledClient struct {
	ID primitive.ObjectID `m:"Database ID"`
	Name string `m:"Name"`
	Phone string `m:"Phone Number"`
	Labels []string `m:"Labels"`
}

func (cl Client) Create(c *gin.Context, req PPA.Client, editor PPA.Editor) (*PPA.Client, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	// ensure OID is unique
	for cl.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}

	labelNames := make([]string, 0)
	for _, labelOID := range req.Labels {
		clientLabel, labelErr := cl.ldb.ViewById(cl.db, ctx, labelOID); if labelErr != nil {
			return nil, PPA.NewAppError(NotFound, "Label Does Not Exist")
		}

		labelNames = append(labelNames, clientLabel.LabelName)
		clientLabel.AppendClient(req.ID)
		updateErr := cl.ldb.Update(cl.db, ctx, clientLabel.ID, clientLabel); if updateErr != nil {
			return nil, PPA.InternalError
		}
	}

	created, err := cl.cdb.Create(cl.db, ctx, &req); if err != nil {
		return nil, err
	}

	created.AppendLog(cl.eventLogger.LogCreated(ctx, cl.eventLogger.GenerateEvent(&logableLabeledClient {
		ID: created.ID,
		Name: created.Name,
		Phone: created.Phone,
		Labels: labelNames,
	}, EventTag), editor))
	cl.cdb.LogEvent(cl.db, ctx, created)
	return created, nil
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

func (cl Client) Delete(c *gin.Context, id string, editor PPA.Editor) error {
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
	cl.deletejobs(ctx, client.Jobs, editor)
	cl.deleteFromLabels(ctx, client.ID, client.Labels)

	client.AppendLog(cl.eventLogger.LogDeleted(ctx, editor))
	cl.cdb.LogEvent(cl.db, ctx, client)
	return cl.cdb.Delete(cl.db, ctx, oid)

}

func (cl Client) Update(c *gin.Context, req Update, id string, editor PPA.Editor) (*PPA.Client, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	// Note that len(nil) is 0
	if len(req.Phone) > 0 {
		fetched, _ := cl.cdb.ViewByPhone(cl.db, ctx, req.Phone)
		if fetched != nil {
			if fetched.ID.Hex() != id {
				return nil, PPA.NewAppError(Conflict, "Phone vnumber already already in use")
			}
		}
	}

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, OidNotFound
	}

	oldDoc, err := cl.cdb.ViewById(cl.db, ctx, oid); if err != nil {
		return nil, err
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

	updated.AppendLog(cl.eventLogger.LogUpdated(ctx,
		cl.eventLogger.GenerateEvent(oldDoc, EventTag),
		cl.eventLogger.GenerateEvent(updated, EventTag),
		editor))
	cl.cdb.LogEvent(cl.db, ctx, updated)

	// update client info on all corresponding jobs
	cl.updateJobs(ctx, updated.Jobs, JobUpdate { Name: updated.Name, Phone: updated.Phone }, editor)
	return updated, nil
}

func (cl Client) Populate(c *gin.Context, unpopClient *PPA.Client) (*PopulatedClient, error) {
	/*if unpopClient.Jobs == nil {
		return nil, nil
	}*/
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()


	jobs := make([]PPA.Job, 0)
	if unpopClient.Jobs != nil {
		res, err := cl.cdb.PopulateJobs(cl.db, ctx, unpopClient.Jobs); if err != nil {
			return nil, err
		}
		jobs = res;
	}

	labels := make([]string, 0)
	if unpopClient.Labels != nil {
		res, lErr := cl.cdb.PopulateLabels(cl.db, ctx, unpopClient.Labels); if lErr != nil {
			return nil, lErr
		}
		labels = res
	}

	return &PopulatedClient {
		ID: unpopClient.ID,
		Name: unpopClient.Name,
		Phone: unpopClient.Phone,
		Jobs: jobs,
		Labels: labels,
		History: unpopClient.History,
	}, nil
}

func (cl Client) PopulateAll(c *gin.Context, unpopClients *[]PPA.Client) (*[]PopulatedClient, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	var popClients = make([]PopulatedClient, 0, len(*unpopClients))

	for _, unpop := range *unpopClients {
		jobs := make([]PPA.Job, 0)
		if unpop.Jobs != nil {
			res, err := cl.cdb.PopulateJobs(cl.db, ctx, unpop.Jobs); if err != nil {
				//return nil, err 
				// skip
			}
			jobs = res;
		}

		labels := make([]string, 0)
		if unpop.Labels != nil {
			res, lErr := cl.cdb.PopulateLabels(cl.db, ctx, unpop.Labels); if lErr != nil {
				//return nil, lErr
				// skip
			}
			labels = res
		}
		popClients = append(popClients, PopulatedClient {
			ID: unpop.ID,
			Name: unpop.Name,
			Phone: unpop.Phone,
			Jobs: jobs,
			Labels: labels,
			History: unpop.History,
		})
	}
	return &popClients, nil
}

func (cl Client) UpdateLabels(c *gin.Context, labels []string, clientID string, editor PPA.Editor) (*PPA.Client, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	clientOID, err := primitive.ObjectIDFromHex(clientID); if err != nil {
		return nil, OidNotFound
	}

	client, clientErr := cl.cdb.ViewById(cl.db, ctx, clientOID); if clientErr != nil {
		return nil, PPA.NewAppError(NotFound, "Client Not Found")
	}
	oldDoc := *client

	labelOIDs := make([]primitive.ObjectID, 0)
	newDocLabels := make([]string, 0)

	for _, label := range labels {
		clientLabel, labelErr := cl.ldb.ViewByLabelName(cl.db, ctx, label); if labelErr != nil {
			return nil, PPA.NewAppError(NotFound, "Label Does Not Exist")
		}

		if !clientLabel.IsDeleted {
			newDocLabels = append(newDocLabels, clientLabel.LabelName)
			clientLabel.AppendClient(clientOID)
			updateErr := cl.ldb.Update(cl.db, ctx, clientLabel.ID, clientLabel); if updateErr != nil {
				return nil, PPA.InternalError
			}
			labelOIDs = append(labelOIDs, clientLabel.ID)
		}
	}
	client.Labels = labelOIDs

	oldDocLabels := make([]string, 0)
	for _, labelOID := range oldDoc.Labels {
		clientLabel, labelErr := cl.ldb.ViewById(cl.db, ctx, labelOID); if labelErr != nil {
			return nil, PPA.NewAppError(NotFound, "Label Does Not Exist")
		}
		if !clientLabel.IsDeleted {
			oldDocLabels = append(oldDocLabels, clientLabel.LabelName)
		}
	}


	logableOldDoc := logableLabeledClient {
		ID: oldDoc.ID,
		Name: oldDoc.Name,
		Phone: oldDoc.Phone,
		Labels: oldDocLabels,
	}

	logableNewDoc := logableLabeledClient {
		ID: client.ID,
		Name: client.Name,
		Phone: client.Phone,
		Labels: newDocLabels,

	}

	client.AppendLog(cl.eventLogger.LogUpdated(ctx,
		cl.eventLogger.GenerateEvent(logableOldDoc, EventTag),
		cl.eventLogger.GenerateEvent(logableNewDoc, EventTag),
		editor))
	cl.cdb.LogEvent(cl.db, ctx, client)
	return client, nil
}

func (cl Client) FetchLabelOIDs(c *gin.Context, labels []string) ([]primitive.ObjectID, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	labelOIDs := make([]primitive.ObjectID, 0)

	for _, label := range labels {
		clientLabel, labelErr := cl.ldb.ViewByLabelName(cl.db, ctx, label); if labelErr != nil {
			return nil, PPA.NewAppError(NotFound, "Label Does Not Exist")
		}
		if !clientLabel.IsDeleted {
			labelOIDs = append(labelOIDs, clientLabel.ID)
		}
	}
	return labelOIDs, nil
}

func (cl Client) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := cl.db.Use(Collection)

	var inserted PPA.Client
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}

func (cl Client) deletejobs(ctx context.Context, oids []primitive.ObjectID, editor PPA.Editor) {
	for _, oid := range oids {
		deleted, _ := cl.jdb.ViewById(cl.db, ctx, oid)
		if deleted != nil {
			deleted.AppendLog(cl.eventLogger.LogDeleted(ctx, editor))
			cl.jdb.LogEvent(cl.db, ctx, deleted)
		}
		err := cl.jdb.Delete(cl.db, ctx, oid); if err != nil {
			// do nothing
		}
	}
}

func (cl Client) deleteFromLabels(ctx context.Context, clientOID primitive.ObjectID, oids []primitive.ObjectID) {
	for _, oid := range oids {
		label, _ := cl.ldb.ViewById(cl.db, ctx, oid)
		label.FindAndRemoveClient(clientOID)
		_ = cl.ldb.PutLabels(cl.db, ctx, label.ID, label.Clients)
	}
}

// How strong of a consistency garuntee do we want????
func (cl Client) updateJobs(ctx context.Context, oids []primitive.ObjectID, update JobUpdate, editor PPA.Editor) {
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
}
