package job

import (
	"PPA"
	//"fmt"
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	NotFound         = http.StatusNotFound
	Conflict         = http.StatusConflict
	Collection       = "jobs"
	ClientCollection = "clients"
	EventTag         = "m"
)

var OidNotFound = PPA.NewAppError(NotFound, "Does not exist")
var JobNotFound = PPA.NewAppError(NotFound, "Does not exist")

type Update struct {
	Title          string
	ClientName     string
	ClientPhone    string
	AssignedWorker primitive.ObjectID
	CarInfo        string
	Notes          string
	StartTime      string
	EndTime        string
	Tag            string
	//Color *PPA.CalendarMeta
	PrimaryColor   string
	SecondaryColor string
}

type ClientUpdate struct {
	Name  string
	Phone string
}

func (j Job) Create(c *gin.Context, req PPA.Job, editor PPA.Editor) (*PPA.Job, error) {
	duration := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	for j.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}

	var workerAssignChanged bool = false
	var cpyTinter *PPA.Tinter
	const matched int = 0
	if bytes.Compare([]byte(req.AssignedWorker.Hex()), []byte(primitive.NilObjectID.Hex())) != matched {
		// update tinter
		tinter, tinterErr := j.tdb.ViewById(j.db, ctx, req.AssignedWorker)
		if tinterErr == nil {
			workerAssignChanged = true
			cpyTinter = tinter
			_ = j.tdb.AssignJobId(j.db, ctx, tinter.JobsCollection, req.ID)
		} else {
			return nil, PPA.NewAppError(NotFound, "Tinter does not exist")
		}
	}
	created, err := j.jdb.Create(j.db, ctx, &req)
	if err != nil {
		return nil, err
	}

	var loggableTinter *PPA.LoggableTinter
	if cpyTinter != nil {
		loggableTinter = cpyTinter.Loggable()
	} else {
		loggableTinter = nil
	}

	loggableJob := created.Loggable(loggableTinter)

	if workerAssignChanged {
		cpyTinter.AppendLog(j.eventLogger.LogAssignedJob(ctx, j.eventLogger.GenerateEvent(loggableJob, EventTag), editor))
		j.tdb.LogEvent(j.db, ctx, cpyTinter)

	}

	created.AppendLog(j.eventLogger.LogCreated(ctx, j.eventLogger.GenerateEvent(loggableJob, EventTag), editor))
	j.jdb.LogEvent(j.db, ctx, created)

	if !j.clientExists(ctx, created.ClientPhone) {
		// create the client
		if _, err := j.createClient(ctx, &PPA.Client{
			ID:    primitive.NilObjectID,
			Name:  created.ClientName,
			Phone: created.ClientPhone,
			Jobs:  []primitive.ObjectID{},
		}, editor); err != nil {
			return nil, err
		}
	}

	// append id to client
	if err := j.attatchJobToClient(ctx, created.ClientPhone, created, editor); err != nil {
		return nil, err
	}

	return created, nil
}

func (j Job) ViewById(c *gin.Context, id string) (*PPA.Job, error) {
	duration := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, OidNotFound
	}
	return j.jdb.ViewById(j.db, ctx, oid)
}

func (j Job) List(c *gin.Context, opts PPA.BulkFetchOptions) (*[]PPA.Job, error) {
	duration := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	return j.jdb.List(j.db, ctx, opts)
}

func (j Job) Delete(c *gin.Context, id string, editor PPA.Editor) error {
	duration := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return OidNotFound
	}

	job, err := j.jdb.ViewById(j.db, ctx, oid)
	if err != nil {
		return JobNotFound
	}

	if err := j.cdb.RemoveJob(j.db, ctx, job.ClientPhone, oid); err != nil {
		return err
	}

	tinter, tinterErr := j.tdb.ViewById(j.db, ctx, job.AssignedWorker)
	if tinterErr == nil {
		_ = j.tdb.RemoveJobId(j.db, ctx, tinter.JobsCollection, job.AssignedWorker)
		tinter.AppendLog(j.eventLogger.LogDeleted(ctx, editor))
		j.tdb.LogEvent(j.db, ctx, tinter)
	} else {
		return PPA.NewAppError(NotFound, "Tinter does not exist")
	}

	job.AppendLog(j.eventLogger.LogDeleted(ctx, editor))
	j.jdb.LogEvent(j.db, ctx, job)
	return j.jdb.Delete(j.db, ctx, oid)
}

func (j Job) Update(c *gin.Context, req Update, id string, editor PPA.Editor) (*PPA.Job, error) {
	duration := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, OidNotFound
	}

	oldJob, err := j.jdb.ViewById(j.db, ctx, oid)
	if err != nil {
		return nil, err
	}

	client, err := j.cdb.ViewByPhone(j.db, ctx, oldJob.ClientPhone)
	if err != nil {
		return nil, err
	}

	var workerAssignChanged bool = false
	var cpyTinter *PPA.Tinter
	const matched int = 0
	if bytes.Compare([]byte(req.AssignedWorker.Hex()), []byte(primitive.NilObjectID.Hex())) != matched {
		workerAssignChanged = true
		// update tinter
		tinter, tinterErr := j.tdb.ViewById(j.db, ctx, req.AssignedWorker)
		if tinterErr == nil {
			cpyTinter = tinter
			// remove old OID
			_ = j.tdb.RemoveJobId(j.db, ctx, tinter.JobsCollection, oldJob.AssignedWorker)
			// insert the replacement
			_ = j.tdb.AssignJobId(j.db, ctx, tinter.JobsCollection, req.AssignedWorker)
		} else {
			return nil, PPA.NewAppError(NotFound, "Tinter does not exist")
		}
	}

	if len(req.ClientPhone) > 0 {
		fetched, _ := j.cdb.ViewByPhone(j.db, ctx, req.ClientPhone)
		if fetched != nil {
			if fetched.ID.Hex() != client.ID.Hex() {
				return nil, PPA.NewAppError(Conflict, "Phone number already already in use")
			}
		}
	}

	if err := j.jdb.Update(j.db, ctx, oid, &PPA.Job{
		ID:             primitive.NilObjectID,
		Title:          req.Title,
		ClientName:     req.ClientName,
		ClientPhone:    req.ClientPhone,
		AssignedWorker: req.AssignedWorker,
		CarInfo:        req.CarInfo,
		Notes:          req.Notes,
		Tag:            req.Tag,
		//Color: req.Color,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
	}); err != nil {
		return nil, err
	}

	if client.Name != req.ClientName || client.Phone != req.ClientPhone {
		if err = j.cdb.Update(j.db, ctx, client.ID, &PPA.Client{
			Name:  req.ClientName,
			Phone: req.ClientPhone,
		}); err != nil {
			return nil, err
		}
		j.updateJobs(ctx, client.Jobs, oid, ClientUpdate{Name: req.ClientName, Phone: req.ClientPhone}, editor)
	}

	updated, err := j.jdb.ViewById(j.db, ctx, oid)
	if err != nil {
		return nil, PPA.InternalError
	}

	var loggableTinter *PPA.LoggableTinter
	if cpyTinter != nil {
		loggableTinter = cpyTinter.Loggable()
	} else {
		loggableTinter = nil
	}
	loggableJob := updated.Loggable(loggableTinter)

	// fetch the old job's tinter for logging
	var loggableOldTinter *PPA.LoggableTinter
	oldTinter, oldTErr := j.tdb.ViewById(j.db, ctx, oldJob.AssignedWorker)
	if oldTErr != nil {
		loggableOldTinter = nil
	} else {
		loggableOldTinter = oldTinter.Loggable()
	}

	loggableOldJob := oldJob.Loggable(loggableOldTinter)

	if workerAssignChanged {
		cpyTinter.AppendLog(j.eventLogger.LogAssignedJob(ctx, j.eventLogger.GenerateEvent(loggableJob, EventTag), editor))
		j.tdb.LogEvent(j.db, ctx, cpyTinter)
	}

	updated.AppendLog(j.eventLogger.LogUpdated(ctx,
		j.eventLogger.GenerateEvent(loggableOldJob, EventTag),
		j.eventLogger.GenerateEvent(loggableJob, EventTag),
		editor))
	j.jdb.LogEvent(j.db, ctx, updated)

	return updated, nil
}

func (j Job) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := j.db.Use(Collection)

	var inserted PPA.Job
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}

func (j Job) clientExists(ctx context.Context, phone string) bool {
	_, err := j.cdb.ViewByPhone(j.db, ctx, phone)
	return err == nil
}

func (j Job) createClient(ctx context.Context, client *PPA.Client, editor PPA.Editor) (*PPA.Client, error) {
	client.ID = primitive.NewObjectID()
	for j.clientOidExists(ctx, client.ID) {
		client.ID = primitive.NewObjectID()
	}

	created, err := j.cdb.Create(j.db, ctx, client)
	if err != nil {
		return nil, err
	}

	created.AppendLog(j.eventLogger.LogCreated(ctx, j.eventLogger.GenerateEvent(created, EventTag), editor))
	j.cdb.LogEvent(j.db, ctx, created)

	return created, nil
}

func (j Job) clientOidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := j.db.Use(ClientCollection)

	var inserted PPA.Job
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}

func (j Job) attatchJobToClient(ctx context.Context, phone string, job *PPA.Job, editor PPA.Editor) error {
	client, err := j.cdb.ViewByPhone(j.db, ctx, phone)
	if err != nil {
		return err
	}

	client.AttatchJob(job.ID)
	updateErr := j.cdb.Update(j.db, ctx, client.ID, client)

	client.AppendLog(j.eventLogger.LogAssignedJob(ctx, j.eventLogger.GenerateEvent(job, EventTag), editor))
	j.cdb.LogEvent(j.db, ctx, client)
	return updateErr
}

func (j Job) updateJobs(ctx context.Context, oids []primitive.ObjectID, currOID primitive.ObjectID, update ClientUpdate, editor PPA.Editor) {
	if len(oids) <= 0 {
		return
	} // Note that len(nil) is 0

	for _, oid := range oids {
		const matched int = 0
		if bytes.Compare([]byte(oid.Hex()), []byte(currOID.Hex())) != matched {
			oldDoc, _ := j.jdb.ViewById(j.db, ctx, oid)

			if err := j.jdb.Update(j.db, ctx, oid, &PPA.Job{
				ClientName:  update.Name,
				ClientPhone: update.Phone,
			}); err != nil {
				// do nothing or fail, depends
			}

			newDoc, _ := j.jdb.ViewById(j.db, ctx, oid)

			if newDoc != nil && oldDoc != nil {
				newDoc.AppendLog(j.eventLogger.LogUpdated(ctx,
					j.eventLogger.GenerateEvent(oldDoc, EventTag),
					j.eventLogger.GenerateEvent(newDoc, EventTag),
					editor))
				j.jdb.LogEvent(j.db, ctx, newDoc)
			}
		}
	}
}
