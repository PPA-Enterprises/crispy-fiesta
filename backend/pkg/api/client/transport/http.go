package transport

import(
	"PPA"
	"strconv"
	"net/http"
	"pkg/api/client"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HTTP struct {
	service client.Service
}

const (
	BadRequest = http.StatusBadRequest
)

func NewHTTP(service client.Service, router *gin.RouterGroup, authMw gin.HandlerFunc) {
	httpTransport := HTTP{service}
	routes := router.Group("/clients")
	routes.POST("/", httpTransport.create)
	routes.GET("/", httpTransport.list)
	routes.GET("/id/:id", httpTransport.viewById)
	routes.GET("/phone/:phone", httpTransport.viewByPhone)
	routes.PATCH("/:id", httpTransport.update)
	routes.DELETE("/:id", httpTransport.delete)
	routes.PUT("/:id", httpTransport.putClientLabels)
}

func (h HTTP) create(c *gin.Context) {
	//check that user is allowed to make this request

	var data createClientRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err); return
	}

	oid := primitive.NewObjectID()
	editor := PPA.Editor {
		OID: oid,
		Name: "Bob",
		Collection: "events" + oid.Hex() + "a",
	}

	newClient, labelErr := h.tryFromCreateClientRequest(c, &data); if labelErr != nil {
		PPA.Response(c, labelErr); return
	}

	created, err := h.service.Create(c, newClient, editor); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusCreated, clientCreated(created)); return
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

	clients, err := h.service.List(c, options); if err != nil {
		PPA.Response(c, err); return
	}
	populated, err := h.service.PopulateJobs(c, clients); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetchedAll(populated)); return
}

func (h HTTP) viewById(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	fetchedClient, err := h.service.ViewById(c, id); if err != nil {
		PPA.Response(c, err); return
	}

	populated, err := h.service.PopulateJob(c, fetchedClient); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetched(populated)); return
}

func (h HTTP) putClientLabels(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	var data putLabelsRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err); return
	}

	oid := primitive.NewObjectID()
	editor := PPA.Editor {
		OID: oid,
		Name: "Bob",
		Collection: "events" + oid.Hex() + "a",
	}

	fetchedClient, err := h.service.UpdateLabels(c, data.Labels, id, editor); if err != nil {
		PPA.Response(c, err); return
	}

	populated, err := h.service.PopulateJob(c, fetchedClient); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetched(populated)); return
}

func (h HTTP) viewByPhone(c *gin.Context) {
	phone := c.Param("phone")
	if len(phone) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Phone Number Required")); return
	}

	fetchedClient, err := h.service.ViewByPhone(c, phone); if err != nil {
		PPA.Response(c, err); return
	}

	populated, err := h.service.PopulateJob(c, fetchedClient); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetched(populated)); return
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

	populated, err := h.service.PopulateJob(c, updated); if err != nil {
		if populated == nil {
			c.JSON(http.StatusOK, clientUpdatedUnpop(updated)); return
		}
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, clientUpdatedUnpop(updated)); return
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
	return gin.H{"success": true, "message": "Client And Associated Jobs Deleted"}
}

func fetched(c *client.PopulatedClient) gin.H {
	return gin.H{"success": true, "message": "Fetched Client", "payload": c}
}

func fetchedAll(c *[]client.PopulatedClient) gin.H {
	return gin.H{"success": true, "message": "Fetched Clients", "payload": c}
}

func clientCreated(c *PPA.Client) gin.H {
	return gin.H{"success": true, "message": "Client Created", "payload": c}
}

func clientUpdated(c *client.PopulatedClient) gin.H {
	return gin.H{"success": true, "message": "Client Updated", "payload": c}
}

func clientUpdatedUnpop(c *PPA.Client) gin.H {
	return gin.H{"success": true, "message": "Client Updated", "payload": c}
}
