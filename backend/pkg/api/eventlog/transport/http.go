
package transport

import(
	"PPA"
	"fmt"
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

	id, err := c.DefaultQuery("id", ""); if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (id) Param")); return
	}

	evlog, err := h.service.List(c, options, id); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetchedAll(evlog)); return
}

func fetchedAll(logs *[]PPA.LogEvent) gin.H {
	return gin.H{"success": true, "message": "Fetched Clients", "payload": logs}
}
