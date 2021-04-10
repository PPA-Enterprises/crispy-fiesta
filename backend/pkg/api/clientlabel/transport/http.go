package transport

import(
	"PPA"
	"net/http"
	"pkg/api/clientlabel"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HTTP struct {
	service clientlabel.Service
}

const (
	BadRequest = http.StatusBadRequest
)

func NewHTTP(service clientlabel.Service, router *gin.RouterGroup, authMw gin.HandlerFunc) {
	httpTransport := HTTP{service}
	routes := router.Group("/clientlabels")
	routes.POST("/", authMw, httpTransport.create)
	routes.GET("/", authMw, httpTransport.list)
	routes.GET("id/:id", authMw, httpTransport.viewById)
	routes.PATCH("/:id", authMw, httpTransport.update)
	routes.DELETE("/:id", authMw, httpTransport.delete)
}

func (h HTTP) create(c *gin.Context) {
	var data createRequest
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

	newLabel := h.fromCreateRequest(&data)
	created, err := h.service.Create(c, newLabel, editor); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusCreated, labelCreated(created)); return
}

func (h HTTP) viewById(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	fetchedLabel, err := h.service.ViewById(c, id); if err != nil {
		PPA.Response(c, err); return
	}

	c.JSON(http.StatusOK, fetched(fetchedLabel)); return
}

func (h HTTP) list (c *gin.Context) {
	labels, err := h.service.List(c); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetchedAll(labels)); return
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
	c.JSON(http.StatusOK, labelUpdated(updated)); return
}


func deleted() gin.H {
	return gin.H{"success": true, "message": "Client Label Deleted"}
}

func fetched(u *PPA.ClientLabel) gin.H {
	return gin.H{"success": true, "message": "Fetched Client Label", "payload": u}
}

func fetchedAll(u *[]PPA.ClientLabel) gin.H {
	return gin.H{"success": true, "message": "Fetched Client Labels", "payload": u}
}

func labelCreated(u *PPA.ClientLabel) gin.H {
	return gin.H{"success": true, "message": "Client Label Created", "payload": u}
}

func labelUpdated(u *PPA.ClientLabel) gin.H {
	return gin.H{"success": true, "message": "Label Updated", "payload": u}
}
