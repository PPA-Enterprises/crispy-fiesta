package clientlabel
import (
	"PPA"
	"pkg/common/mongo"
	"github.com/gin-gonic/gin"
	"context"
	dbQuery "pkg/api/clientlabel/infra/mongo"
)

type ClientLabelStream struct {
	db *mongo.DBConnection
	cldb StreamRepository
	rbac RBAC
}


type StreamRepository interface {
	Stream(*mongo.DBConnection, context.Context, chan *PPA.StreamResult)
}

type StreamService interface {
	Subscribe(*gin.Context, *PPA.StreamEvent) error
}

func NewStream(db *mongo.DBConnection, cldb StreamRepository, rbac RBAC) *ClientLabelStream {
	return &ClientLabelStream{db: db, cldb: cldb, rbac: rbac}
}

func InitStream(db *mongo.DBConnection, rbac RBAC) *ClientLabelStream {
	return NewStream(db, dbQuery.ClientLabel{}, rbac)
}
