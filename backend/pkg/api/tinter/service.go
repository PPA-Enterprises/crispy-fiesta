package tinter

import (
	"PPA"
	"context"
	"pkg/common/mongo"
	dbQuery "pkg/api/tinter/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type Tinter struct {
	db *mongo.DBConnection
	tdb Repository
	jdb JobRepository
	eventLogger EventLogger
	rbac RBAC
}

type Service interface {
	Create(*gin.Context, PPA.Tinter, PPA.Editor) (*PPA.Tinter, error)
	List(*gin.Context, PPA.BulkFetchOptions) (*[]PPA.Tinter, error)
	ViewById(*gin.Context, string) (*PPA.Tinter, error)
	ViewByPhone(*gin.Context, string) (*PPA.Tinter, error)
	Delete(*gin.Context, string, PPA.Editor) error
	Update(*gin.Context, Update, string, PPA.Editor) (*PPA.Tinter, error)
	//PopulateJob(*gin.Context, *PPA.Tinter) (*PopulatedTinter, error)
	//PopulateJobs(*gin.Context, *[]PPA.Tinter) (*[]PopulatedTinter, error)
}

func New(db *mongo.DBConnection, tdb Repository, jdb JobRepository, rbac RBAC, ev EventLogger) *Tinter {
	return &Tinter{db: db, tdb: tdb, jdb: jdb, eventLogger: ev, rbac: rbac}
}

func Init(db *mongo.DBConnection, rbac RBAC, jdb JobRepository, ev EventLogger) *Tinter {
	return New(db, dbQuery.Tinter{}, jdb, rbac, ev)
}

type JobRepository interface {
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Job) error
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Job, error)
	LogEvent(*mongo.DBConnection, context.Context, *PPA.Job)
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.Tinter) (*PPA.Tinter, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Tinter, error)
	ViewByPhone(*mongo.DBConnection, context.Context, string) (*PPA.Tinter, error)
	List(*mongo.DBConnection, context.Context, PPA.BulkFetchOptions) (*[]PPA.Tinter, error)
	Delete(*mongo.DBConnection, context.Context, primitive.ObjectID) error
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Tinter) error
	PutJobs(*mongo.DBConnection, context.Context, primitive.ObjectID, []primitive.ObjectID) error
	Populate(*mongo.DBConnection, context.Context, []primitive.ObjectID) ([]PPA.Job, error)
	LogEvent(*mongo.DBConnection, context.Context, *PPA.Tinter)
}

type EventLogger interface {
	LogCreated(context.Context, PPA.EventMap, PPA.Editor) PPA.LogEvent
	LogUpdated(context.Context, PPA.EventMap, PPA.EventMap, PPA.Editor) PPA.LogEvent
	LogDeleted(context.Context, PPA.Editor) PPA.LogEvent
	GenerateEvent(interface{}, string) PPA.EventMap
}
/*
type PopulatedTinter struct {
	ID primitive.ObjectID `json:"_id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Jobs []PPA.Job `json:"jobs"`
	History []PPA.LogEvent `json:"history"`
}*/

type RBAC interface {
}
