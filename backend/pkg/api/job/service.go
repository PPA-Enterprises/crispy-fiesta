package job
import(
	"PPA"
	"context"
	"pkg/common/mongo"
	dbQuery "pkg/api/job/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type Job struct {
	db *mongo.DBConnection
	jdb Repository
	cdb ClientRepository
	tdb TinterRepository
	eventLogger EventLogger
	rbac RBAC
}

type Service interface {
	Create(*gin.Context, PPA.Job, PPA.Editor) (*PPA.Job, error)
	List(*gin.Context, PPA.BulkFetchOptions) (*[]PPA.Job, error)
	ViewById(*gin.Context, string) (*PPA.Job, error)
	Delete(*gin.Context, string, PPA.Editor) error
	Update(*gin.Context, Update, string, PPA.Editor) (*PPA.Job, error)
}

func New(db *mongo.DBConnection, jdb Repository, cdb ClientRepository, tdb TinterRepository, ev EventLogger, rbac RBAC) *Job {
	return &Job{db: db, jdb: jdb, cdb: cdb, tdb: tdb, eventLogger: ev, rbac: rbac}
}

func Init(db *mongo.DBConnection, rbac RBAC, cdb ClientRepository, tdb TinterRepository, ev EventLogger) *Job {
	return New(db, dbQuery.Job{}, cdb, tdb, ev, rbac)
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.Job) (*PPA.Job, error)
	List(*mongo.DBConnection, context.Context, PPA.BulkFetchOptions) (*[]PPA.Job, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Job, error)
	Delete(*mongo.DBConnection, context.Context, primitive.ObjectID) error
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Job) error
	LogEvent(*mongo.DBConnection, context.Context, *PPA.Job)
}

type ClientRepository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.Client) (*PPA.Client, error)
	ViewByPhone(*mongo.DBConnection, context.Context, string) (*PPA.Client, error)
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Client) error
	RemoveJob(*mongo.DBConnection, context.Context, string, primitive.ObjectID) error
	LogEvent(*mongo.DBConnection, context.Context, *PPA.Client)
}

type TinterRepository interface {
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Tinter, error)
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Tinter) error
	PutJobs(*mongo.DBConnection, context.Context, primitive.ObjectID, []primitive.ObjectID) error
}

type EventLogger interface {
	LogCreated(context.Context, PPA.EventMap, PPA.Editor) PPA.LogEvent
	LogUpdated(context.Context, PPA.EventMap, PPA.EventMap, PPA.Editor) PPA.LogEvent
	LogDeleted(context.Context, PPA.Editor) PPA.LogEvent
	LogAssignedJob(context.Context, PPA.EventMap, PPA.Editor) PPA.LogEvent
	GenerateEvent(interface{}, string) PPA.EventMap
}

type RBAC interface{}
