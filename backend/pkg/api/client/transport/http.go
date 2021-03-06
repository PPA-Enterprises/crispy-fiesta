package transport

import(
	"PPA"
	"fmt"
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
	routes.POST("/", authMw, httpTransport.create)
	routes.GET("/", authMw, httpTransport.list)
	routes.GET("/id/:id", authMw, httpTransport.viewById)
	routes.GET("/phone/:phone", authMw, httpTransport.viewByPhone)
	routes.PATCH("/:id", authMw, httpTransport.update)
	routes.DELETE("/:id", authMw, httpTransport.delete)
	routes.PUT("/labels/:id", authMw, httpTransport.putClientLabels)
}

func (h HTTP) create(c *gin.Context) {
	//check that user is allowed to make this request

	var data createClientRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err); return
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
	populated, err := h.service.PopulateAll(c, clients); if err != nil {
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

	populated, err := h.service.Populate(c, fetchedClient); if err != nil {
		PPA.Response(c, err); return
	}
	fmt.Println(populated);
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

	fetchedClient, err := h.service.UpdateLabels(c, data.Labels, id, editor); if err != nil {
		PPA.Response(c, err); return
	}

	populated, err := h.service.Populate(c, fetchedClient); if err != nil {
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

	populated, err := h.service.Populate(c, fetchedClient); if err != nil {
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

	updated, err := h.service.Update(c, h.fromUpdateRequest(&data), id, editor); if err != nil {
		PPA.Response(c, err); return
	}

	populated, err := h.service.Populate(c, updated); if err != nil {
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
