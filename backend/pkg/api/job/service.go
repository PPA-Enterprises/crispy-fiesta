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
	rbac RBAC
}

type Service interface {
	Create(*gin.Context, PPA.Job) (*PPA.Job, error)
	List(*gin.Context, PPA.BulkFetchOptions) (*[]PPA.Job, error)
	ViewById(*gin.Context, string) (*PPA.Job, error)
	Delete(*gin.Context, string) error
	Update(*gin.Context, Update, string) (*PPA.Job, error)
}

func New(db *mongo.DBConnection, jdb Repository, cdb ClientRepository, rbac RBAC) *Job {
	return &Job{db: db, jdb: jdb, cdb: cdb, rbac: rbac}
}

func Init(db *mongo.DBConnection, rbac RBAC, cdb ClientRepository) *Job {
	return New(db, dbQuery.Job{}, cdb, rbac)
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.Job) (*PPA.Job, error)
	List(*mongo.DBConnection, context.Context, PPA.BulkFetchOptions) (*[]PPA.Job, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Job, error)
	Delete(*mongo.DBConnection, context.Context, primitive.ObjectID) error
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Job) error
}

type ClientRepository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.Client) (*PPA.Client, error)
	ViewByPhone(*mongo.DBConnection, context.Context, string) (*PPA.Client, error)
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Client) error
	RemoveJob(*mongo.DBConnection, context.Context, string, primitive.ObjectID) error
}

type RBAC interface{}
