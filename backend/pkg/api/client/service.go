package client

import (
	"PPA"
	"context"
	"pkg/common/mongo"
	dbQuery "pkg/api/client/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type Client struct {
	db *mongo.DBConnection
	cdb Repository
	jdb JobRepository
	eventLogger EventLogger
	rbac RBAC
}

type Service interface {
	Create(*gin.Context, PPA.Client, PPA.Editor) (*PPA.Client, error)
	List(*gin.Context, PPA.BulkFetchOptions) (*[]PPA.Client, error)
	ViewById(*gin.Context, string) (*PPA.Client, error)
	ViewByPhone(*gin.Context, string) (*PPA.Client, error)
	Delete(*gin.Context, string, PPA.Editor) error
	Update(*gin.Context, Update, string, PPA.Editor) (*PPA.Client, error)
	PopulateJob(*gin.Context, *PPA.Client) (*PopulatedClient, error)
	PopulateJobs(*gin.Context, *[]PPA.Client) (*[]PopulatedClient, error)
}

func New(db *mongo.DBConnection, cdb Repository, jdb JobRepository, rbac RBAC, ev EventLogger) *Client {
	return &Client{db: db, cdb: cdb, jdb: jdb, eventLogger: ev, rbac: rbac}
}

func Init(db *mongo.DBConnection, rbac RBAC, jdb JobRepository, ev EventLogger) *Client {
	return New(db, dbQuery.Client{}, jdb, rbac, ev)
}

type JobRepository interface {
	Delete(*mongo.DBConnection, context.Context, primitive.ObjectID) error
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Job) error
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.Client) (*PPA.Client, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Client, error)
	ViewByPhone(*mongo.DBConnection, context.Context, string) (*PPA.Client, error)
	List(*mongo.DBConnection, context.Context, PPA.BulkFetchOptions) (*[]PPA.Client, error)
	Delete(*mongo.DBConnection, context.Context, primitive.ObjectID) error
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Client) error
	Populate(*mongo.DBConnection, context.Context, []primitive.ObjectID) ([]PPA.Job, error)
	LogEvent(*mongo.DBConnection, context.Context, *PPA.Client)
}

type EventLogger interface {
	LogCreated(context.Context, PPA.EventMap, PPA.Editor) PPA.LogEvent
	LogUpdated(context.Context, PPA.EventMap, PPA.EventMap, PPA.Editor) PPA.LogEvent
	LogDeleted(context.Context, PPA.Editor) PPA.LogEvent
	GenerateEvent(interface{}, string) PPA.EventMap
}

type PopulatedClient struct {
	ID primitive.ObjectID `json:"_id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Jobs []PPA.Job `json:"jobs"`
}

type RBAC interface {
}
