package eventlog
import (
	"PPA"
	"context"
	"pkg/common/mongo"
	dbQuery "pkg/api/eventlog/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type Eventlog struct {
	db *mongo.DBConnection
	evdb Repository,
	rbac RBAC
}

type Service interface {
	List(*gin.Context, PPA.BulkFetchOptions) (*[]PPA.LogEvent, error)
}

func New(db *mongo.DBConnection, evdb Repository, rbac RBAC) *Eventlog {
	return &Eventlog{db: db, evdb: evdb, rbac: rbac}
}

func Init(db *mongo.DBConnection, rbac RBAC) *Eventlog {
	return New(db, dbQuery.Eventlog{}, rbac)
}

type Repository interface {
	List(*gin.Context, context.Context, PPA.BulkFetchOptions) (*[]PPA.LogEvent, error)
}

type RBAC interface{}
