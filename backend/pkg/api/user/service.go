package user

import(
	"PPA"
	"context"
	"pkg/common/mongo"
	dbQuery "pkg/api/user/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(*gin.Context, PPA.User, PPA.Editor) (*PPA.User, error)
	List(*gin.Context) (*[]PPA.User, error)
	ViewById(*gin.Context, string) (*PPA.User, error)
	ViewByEmail(*gin.Context, string) (*PPA.User, error)
	Delete(*gin.Context, string, PPA.Editor) error
	Update(*gin.Context, Update, string, PPA.Editor) (*PPA.User, error)
	ListEvents(*gin.Context, PPA.BulkFetchOptions, string) (*[]PPA.LogEvent, error)
}

func New(db *mongo.DBConnection, udb Repository, rbac RBAC, securer Securer, el EventLogger) *User {
	return &User{db: db, udb: udb, rbac: rbac, securer: securer, eventLogger: el}
}

func Init(db *mongo.DBConnection, rbac RBAC, securer Securer, el EventLogger) *User {
	return New(db, dbQuery.User{}, rbac, securer, el)
}

type User struct {
	db *mongo.DBConnection
	udb Repository
	rbac RBAC
	securer Securer
	eventLogger EventLogger
}

type Securer interface {
	Hash(string) string
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.User) (*PPA.User, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.User, error)
	ViewByEmail(*mongo.DBConnection, context.Context, string) (*PPA.User, error)
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.User) error
	List(*mongo.DBConnection, context.Context) (*[]PPA.User, error)
	Delete(*mongo.DBConnection, context.Context, primitive.ObjectID) error
	LogEvent(*mongo.DBConnection, context.Context, *PPA.User)
	ListEvents(*mongo.DBConnection, context.Context, PPA.BulkFetchOptions, string) (*[]PPA.LogEvent, error)
}

type EventLogger interface {
	LogCreated(context.Context, PPA.EventMap, PPA.Editor) PPA.LogEvent
	LogUpdated(context.Context, PPA.EventMap, PPA.EventMap, PPA.Editor) PPA.LogEvent
	LogDeleted(context.Context, PPA.Editor) PPA.LogEvent
	GenerateEvent(interface{}, string) PPA.EventMap
}

type RBAC interface {
	//User(*gin.Context) PPA.AuthUser
	//AccountCreate(*gin.Context) error
}
