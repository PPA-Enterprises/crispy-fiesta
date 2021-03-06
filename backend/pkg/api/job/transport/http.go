package transport

import (
	"PPA"
	"net/http"
	"pkg/api/job"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HTTP struct {
	service job.Service
}

const (
	BadRequest = http.StatusBadRequest
)

func NewHTTP(service job.Service, router *gin.RouterGroup, authMw gin.HandlerFunc) {
	httpTransport := HTTP{service}
	routes := router.Group("/jobs")
	routes.POST("/", authMw, httpTransport.create)
	routes.GET("/", authMw, httpTransport.list)
	routes.GET("id/:id", authMw, httpTransport.viewById)
	routes.PATCH("/:id", authMw, httpTransport.update)
	routes.DELETE("/:id", authMw, httpTransport.delete)
}

func (h HTTP) create(c *gin.Context) {
	//check that user is allowed to make this request

	var data submitJobRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err)
		return
	}

	newJob, reqErr := h.tryFromSubmitJobRequest(&data); if reqErr != nil {
		PPA.Response(c, reqErr)
		return
	}
	editorId := c.GetString(PPA.IdKey)
	editorName := c.GetString(PPA.NameKey)

	editorOid, err := primitive.ObjectIDFromHex(editorId); if err != nil {
		PPA.Response(c, err); return
	}

	editor := PPA.Editor {
		OID: editorOid,
		Name: editorName,
		Collection: "events" + editorOid.Hex() + "a",
	}

	created, err := h.service.Create(c, newJob, editor)
	if err != nil {
		PPA.Response(c, err)
		return
	}
	c.JSON(http.StatusCreated, jobCreated(created))
	return
}

func (h HTTP) viewById(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required"))
		return
	}

	fetchedJob, err := h.service.ViewById(c, id)
	if err != nil {
		PPA.Response(c, err)
		return
	}

	c.JSON(http.StatusOK, fetched(fetchedJob))
	return
}

func (h HTTP) list(c *gin.Context) {
	options := PPA.NewBulkFetchOptions()

	all, err := strconv.ParseBool(c.DefaultQuery("all", "false"))
	if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (all) Param"))
		return
	}
	options.All = all

	sort, err := strconv.ParseBool(c.DefaultQuery("sort", "false"))
	if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (sort) Param"))
		return
	}
	options.Sort = sort

	source, err := strconv.ParseUint(c.DefaultQuery("source", "0"), 10, 64)
	if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (source) Param"))
		return
	}
	options.Source = source

	next, err := strconv.ParseUint(c.DefaultQuery("next", "10"), 10, 64)
	if err != nil {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Invalid (next) Param"))
		return
	}
	options.Next = next

	jobs, err := h.service.List(c, options)
	if err != nil {
		PPA.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, fetchedAll(jobs))
	return
}

func (h HTTP) delete(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required"))
		return
	}

	editorId := c.GetString(PPA.IdKey)
	editorName := c.GetString(PPA.NameKey)

	editorOid, err := primitive.ObjectIDFromHex(editorId); if err != nil {
		PPA.Response(c, err); return
	}

	editor := PPA.Editor {
		OID: editorOid,
		Name: editorName,
		Collection: "events" + editorOid.Hex() + "a",
	}

	if err := h.service.Delete(c, id, editor); err != nil {
		PPA.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, deleted())
	return
}

func (h HTTP) update(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required"))
		return
	}

	var data updateRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err)
		return
	}

	editorId := c.GetString(PPA.IdKey)
	editorName := c.GetString(PPA.NameKey)

	editorOid, err := primitive.ObjectIDFromHex(editorId); if err != nil {
		PPA.Response(c, err); return
	}

	editor := PPA.Editor {
		OID: editorOid,
		Name: editorName,
		Collection: "events" + editorOid.Hex() + "a",
	}

	updateData, reqErr := h.tryFromUpdateRequest(&data); if reqErr != nil {
		PPA.Response(c, reqErr)
		return
	}
	updated, err := h.service.Update(c, updateData, id, editor)
	if err != nil {
		PPA.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, jobUpdated(updated))
	return
}

func deleted() gin.H {
	return gin.H{"success": true, "message": "Job Deleted"}
}

func fetched(j *PPA.Job) gin.H {
	return gin.H{"success": true, "message": "Fetched Job", "payload": j}
}

func fetchedAll(j *[]PPA.Job) gin.H {
	return gin.H{"success": true, "message": "Fetched Job", "payload": j}
}

func jobCreated(j *PPA.Job) gin.H {
	return gin.H{"success": true, "message": "Job Created", "payload": j}
}

func jobUpdated(j *PPA.Job) gin.H {
	return gin.H{"success": true, "message": "Job Updated", "payload": j}
}
