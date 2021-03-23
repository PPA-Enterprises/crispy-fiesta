package clientlabel

import(
	"PPA"
	"context"
	"pkg/common/mongo"
	dbQuery "pkg/api/clientlabel/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(*gin.Context, PPA.ClientLabel, PPA.Editor) (*PPA.ClientLabel, error)
	List(*gin.Context) (*[]PPA.ClientLabel, error)
	ViewById(*gin.Context, string) (*PPA.ClientLabel, error)
	Delete(*gin.Context, string, PPA.Editor) error
	Update(*gin.Context, Update, string, PPA.Editor) (*PPA.ClientLabel, error)
}

func New(db *mongo.DBConnection, cldb Repository, rbac RBAC, el EventLogger) *ClientLabel {
	return &ClientLabel{db: db, cldb: cldb, rbac: rbac, eventLogger: el}
}

func Init(db *mongo.DBConnection, rbac RBAC, el EventLogger) *ClientLabel {
	return New(db, dbQuery.ClientLabel{}, rbac, el)
}

type ClientLabel struct {
	db *mongo.DBConnection
	cldb Repository
	rbac RBAC
	eventLogger EventLogger
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.ClientLabel) (*PPA.ClientLabel, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.ClientLabel, error)
	ViewByLabelName(*mongo.DBConnection, context.Context, string) (*PPA.ClientLabel, error)
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.ClientLabel) error
	List(*mongo.DBConnection, context.Context) (*[]PPA.ClientLabel, error)
	LogEvent(*mongo.DBConnection, context.Context, *PPA.ClientLabel)
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
