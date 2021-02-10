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
	evdb Repository
	udb UserRepository
	rbac RBAC
}

type Service interface {
	List(*gin.Context, PPA.BulkFetchOptions, string) (*[]PPA.LogEvent, error)
}

func New(db *mongo.DBConnection, evdb Repository, udb UserRepository, rbac RBAC) *Eventlog {
	return &Eventlog{db: db, evdb: evdb, udb: udb, rbac: rbac}
}

func Init(db *mongo.DBConnection, udb UserRepository, rbac RBAC) *Eventlog {
	return New(db, dbQuery.Eventlog{}, udb, rbac)
}

type Repository interface {
	List(*mongo.DBConnection, context.Context, PPA.BulkFetchOptions, string) (*[]PPA.LogEvent, error)
}

type UserRepository interface {
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.User, error)
}

type RBAC interface{}
