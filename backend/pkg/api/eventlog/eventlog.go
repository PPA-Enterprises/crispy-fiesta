package eventlog
import(
	"PPA"
	"fmt"
	"context"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

var OidNotFound = PPA.NewAppError(http.StatusNotFound, "Does not exist")

func (ev Eventlog) List(c *gin.Context, opts PPA.BulkFetchOptions, id string) (*[]PPA.LogEvent, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return OidNotFound
	}

	user, err := ev.udb.ViewById(ev.db, ctx, oid); if err != nil {
		return OidNotFound
	}

	collection := user.Name + user.ID.Hex() + "a"
	return ev.evdb.List(ev.db, ctx, opts, collection)
}
