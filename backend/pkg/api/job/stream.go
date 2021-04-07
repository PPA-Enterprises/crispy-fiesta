package job
import (
	"PPA"
	"pkg/common/mongo"
	"github.com/gin-gonic/gin"
	"context"
	dbQuery "pkg/api/job/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobStream struct {
	db *mongo.DBConnection
	tdb TinterViewer
	jdb StreamRepository
	rbac RBAC
}


type StreamRepository interface {
	Stream(*mongo.DBConnection, context.Context, chan *PPA.Job)
}

type StreamService interface {
	Subscribe(*gin.Context, *PPA.StreamEvent) error
}

type TinterViewer interface {
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Tinter, error)
}

func NewStream(db *mongo.DBConnection, jdb StreamRepository, tdb TinterViewer, rbac RBAC) *JobStream {
	return &JobStream{db: db, jdb: jdb, tdb: tdb, rbac: rbac}
}

func InitStream(db *mongo.DBConnection, rbac RBAC, tdb TinterViewer) *JobStream {
	return NewStream(db, dbQuery.Job{}, tdb, rbac)
}
