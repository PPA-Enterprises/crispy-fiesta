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
	rbac RBAC
}

type Service interface {
	Create(*gin.Context, PPA.Client) (*PPA.Client, error)
	List(*gin.Context) (*[]PPA.Client, error)
	ViewById(*gin.Context, string) (*PPA.Client, error)
	ViewByPhone(*gin.Context, string) (*PPA.Client, error)
	Delete(*gin.Context, string) error
	Update(*gin.Context, Update, string) (*PPA.Client, error)
	PopulateJob(*gin.Context, *PPA.Client) (*PopulatedClient, error)
//	PopulateJobs(*gin.Context, *[]PPA.Client) (*[]PopulatedClient, error)
}

func New(db *mongo.DBConnection, cdb Repository, rbac RBAC) *Client {
	return &Client{db: db, cdb: cdb, rbac: rbac}
}

func Init(db *mongo.DBConnection, rbac RBAC) *Client {
	return New(db, dbQuery.Client{}, rbac)
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.Client) (*PPA.Client, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Client, error)
	ViewByPhone(*mongo.DBConnection, context.Context, string) (*PPA.Client, error)
	List(*mongo.DBConnection, context.Context) (*[]PPA.Client, error)
	Delete(*mongo.DBConnection, context.Context, primitive.ObjectID) error
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Client) error
	Populate(*mongo.DBConnection, context.Context, []primitive.ObjectID) ([]PPA.Job, error)
}

type PopulatedClient struct {
	ID primitive.ObjectID `json:"_id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Jobs []PPA.Job `json:"jobs"`
}

type RBAC interface {
}
