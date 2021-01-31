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
	rbac RBAC
}

type Service interface {
	Create(*gin.Context, PPA.Job) (*PPA.Job, error)
	List(*gin.Context) (*[]PPA.Job, error)
	ViewById(*gin.Context, string) (*PPA.Job, error)
	Delete(*gin.Context, string) error
	Update(*gin.Context, Update, string) (*PPA.Client, error)
}

func New(db *mongo.DBConnection, jdb Repository, rbac RBAC) *Job {
	return &Job{db: db, jdb: jdb, rbac: rbac}
}

func Init(db *mongo.DBConnection, rbac RBAC) *Job {
	return New(db, dbQuery.Job{}, rbac)
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.Job) (*PPA.Job, error)
	List(*mongo.DBConnection, context.Context) (*[]PPA.Job, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Job, error)
	Delete(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Job, error)
	Update(*mongo.DBConnection, context.Context, primitive.ObjectID, *PPA.Job) error
}

type RBAC interface{}
