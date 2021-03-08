package eventlog
import(
	"PPA"
	"context"
	//"net/http"
	"time"
	"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

//var OidNotFound = PPA.NewAppError(http.StatusNotFound, "Does not exist")

func (ev Eventlog) List(c *gin.Context, opts PPA.BulkFetchOptions) (*[]PPA.LogEvent, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	return ev.evdb.List(ev.db, ctx, opts)
}
