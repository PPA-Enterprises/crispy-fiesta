
package transport

import(
	"PPA"
	"strconv"
	"net/http"
	"pkg/api/eventlog"
	"github.com/gin-gonic/gin"
)

type HTTP struct {
	service eventlog.Service
}

const (
	BadRequest = http.StatusBadRequest
)

func NewHTTP(service eventlog.Service, router *gin.RouterGroup, authMw gin.HandlerFunc) {
	httpTransport := HTTP{service}
	routes := router.Group("/eventlogs")
	routes.GET("/", httpTransport.list)
}

func (h HTTP) list(c *gin.Context) {
	options := PPA.NewBulkFetchOptions()

	all, err := strconv.ParseBool(c.DefaultQuery("all", "false")); if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (all) Param")); return
	}
	options.All = all

	sort, err := strconv.ParseBool(c.DefaultQuery("sort", "false")); if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (sort) Param")); return
	}
	options.Sort = sort

	source, err := strconv.ParseUint(c.DefaultQuery("source", "0"), 10, 64); if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (source) Param")); return
	}
	options.Source = source

	next, err := strconv.ParseUint(c.DefaultQuery("next", "10"), 10, 64); if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (next) Param")); return
	}
	options.Next = next

	evlog, err := h.service.List(c, options); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetchedAll(evlog)); return
}

func fetchedAll(logs *[]PPA.LogEvent) gin.H {
	if logs == nil || len(*logs) <= 0 {
		// JSON serialization bug requires me to do this
		return gin.H{"success": true, "message": "Fetched Logs", "payload": "[]"}
	}
	return gin.H{"success": true, "message": "Fetched Logs", "payload": logs}
}
