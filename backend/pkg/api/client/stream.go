package client
import (
	"PPA"
	"pkg/common/mongo"
	"github.com/gin-gonic/gin"
	"context"
	dbQuery "pkg/api/client/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientStream struct {
	db *mongo.DBConnection
	jdb JobViewer
	ldb LabelViewer
	cdb StreamRepository
	rbac RBAC
}


type StreamRepository interface {
	Stream(*mongo.DBConnection, context.Context, chan *PPA.StreamResult)
	PopulateJobs(*mongo.DBConnection, context.Context, []primitive.ObjectID) ([]PPA.Job, error)
	PopulateLabels(*mongo.DBConnection, context.Context, []primitive.ObjectID) ([]string, error)
}

type StreamService interface {
	Subscribe(*gin.Context, *PPA.StreamEvent) error
}

type JobViewer interface {
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.Job, error)
}

type LabelViewer interface {
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.ClientLabel, error)
	ViewByLabelName(*mongo.DBConnection, context.Context, string) (*PPA.ClientLabel, error)
}

func NewStream(db *mongo.DBConnection, cdb StreamRepository, jdb JobViewer, ldb LabelViewer, rbac RBAC) *ClientStream {
	return &ClientStream{db: db, jdb: jdb, cdb: cdb, ldb:ldb, rbac: rbac}
}

func InitStream(db *mongo.DBConnection, rbac RBAC, jdb JobViewer, ldb LabelViewer) *ClientStream {
	return NewStream(db, dbQuery.Client{}, jdb, ldb, rbac)
}
