package transport

import(
	"PPA"
	"fmt"
	"strconv"
	"net/http"
	"pkg/api/tinter"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HTTP struct {
	service tinter.Service
}

const (
	BadRequest = http.StatusBadRequest
)

func NewHTTP(service tinter.Service, router *gin.RouterGroup, authMw gin.HandlerFunc) {
	httpTransport := HTTP{service}
	routes := router.Group("/tinters")
	routes.POST("/", httpTransport.create)
	routes.GET("/", httpTransport.list)
	routes.GET("/id/:id", httpTransport.viewById)
	routes.GET("/phone/:phone", httpTransport.viewByPhone)
	routes.PATCH("/:id", httpTransport.update)
	routes.DELETE("/:id", httpTransport.delete)
}

func (h HTTP) create(c *gin.Context) {
	//check that user is allowed to make this request

	var data createTinterRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err); return
	}

	oid := primitive.NewObjectID()
	editor := PPA.Editor {
		OID: oid,
		Name: "Bob",
		Collection: "events" + oid.Hex() + "a",
	}

	created, err := h.service.Create(c, h.fromCreateTinterRequest(&data), editor); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusCreated, tinterCreated(created)); return
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

	tinters, err := h.service.List(c, options); if err != nil {
		PPA.Response(c, err); return
	}
	/*populated, err := h.service.PopulateJobs(c, tinters); if err != nil {
		PPA.Response(c, err); return
	}*/
	c.JSON(http.StatusOK, fetchedAll(tinters)); return
}

func (h HTTP) viewById(c *gin.Context) {
	fmt.Println(h)
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	fetchedTinter, err := h.service.ViewById(c, id); if err != nil {
		PPA.Response(c, err); return
	}

	/*populated, err := h.service.PopulateJob(c, fetchedTinter); if err != nil {
		PPA.Response(c, err); return
	}*/
	c.JSON(http.StatusOK, fetched(fetchedTinter)); return
}

func (h HTTP) viewByPhone(c *gin.Context) {
	phone := c.Param("phone")
	if len(phone) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Phone Number Required")); return
	}

	fetchedTinter, err := h.service.ViewByPhone(c, phone); if err != nil {
		PPA.Response(c, err); return
	}
/*
	populated, err := h.service.PopulateJob(c, fetchedTinter); if err != nil {
		PPA.Response(c, err); return
	}*/
	c.JSON(http.StatusOK, fetched(fetchedTinter)); return
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

	oid := primitive.NewObjectID()
	editor := PPA.Editor {
		OID: oid,
		Name: "Bob",
		Collection: "events" + oid.Hex() + "a",
	}

	updated, err := h.service.Update(c, h.fromUpdateRequest(&data), id, editor); if err != nil {
		PPA.Response(c, err); return
	}

	/*populated, err := h.service.PopulateJob(c, updated); if err != nil {
		PPA.Response(c, err); return
	}*/
	c.JSON(http.StatusOK, fetched(updated)); return
}

func (h HTTP) delete(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	oid := primitive.NewObjectID()
	editor := PPA.Editor {
		OID: oid,
		Name: "Bob",
		Collection: "events" + oid.Hex() + "a",
	}

	if err := h.service.Delete(c, id, editor); err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, deleted()); return
}

func deleted() gin.H {
	return gin.H{"success": true, "message": "Tinter Deleted"}
}

func fetched(c *PPA.Tinter) gin.H {
	return gin.H{"success": true, "message": "Fetched Tinter", "payload": c}
}

func fetchedAll(c *[]PPA.Tinter) gin.H {
	return gin.H{"success": true, "message": "Fetched Tinters", "payload": c}
}

func tinterCreated(c *PPA.Tinter) gin.H {
	return gin.H{"success": true, "message": "Tinter Created", "payload": c}
}

func tinterUpdated(c *PPA.Tinter) gin.H {
	return gin.H{"success": true, "message": "Tinter Updated", "payload": c}
}
