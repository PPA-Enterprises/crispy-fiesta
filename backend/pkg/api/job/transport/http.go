package transport

import(
	"PPA"
	"net/http"
	"pkg/api/job"
	"github.com/gin-gonic/gin"
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
	routes.POST("/", httpTransport.create)
	routes.GET("/", httpTransport.list)
	routes.GET("id/:id", httpTransport.viewById)
	routes.PATCH("/:id", httpTransport.update)
	routes.DELETE("/:id", httpTransport.delete)
}

func (h HTTP) create(c *gin.Context) {
	//check that user is allowed to make this request

	var data signupRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err); return
	}

	newUser := h.fromSignupRequest(&data)
	created, err := h.service.Create(c, newUser); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusCreated, jobCreated(created)); return
}

func (h HTTP) viewById(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	fetchedJob, err := h.service.ViewById(c, id); if err != nil {
		PPA.Response(c, err); return
	}

	c.JSON(http.StatusOK, fetched(fetchedJob)); return
}

func (h HTTP) list (c *gin.Context) {
	jobs, err := h.service.List(c); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetchedAll(jobs)); return
}

func (h HTTP) delete(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	if err := h.service.Delete(c, id); err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, deleted()); return
}

func (h HTTP) update(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	var data updateRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err); return
	}

	updated, err := h.service.Update(c, h.fromUpdateRequest(&data), id); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, jobUpdated(updated)); return
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
