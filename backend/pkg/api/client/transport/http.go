package transport

import(
	"PPA"
	"net/http"
	"pkg/api/client"
	"github.com/gin-gonic/gin"
)

type HTTP struct {
	service client.Service
}

func NewHTTP(service client.Service, router *gin.RouterGroup, authMw gin.HandlerFunc) {
	httpTransport := HTTP{service}
	routes := router.Group("/clients")
	routes.POST("/", httpTransport.create)
	routes.GET("/", httpTransport.list)
	routes.GET("/id/:id", httpTransport.viewById)
	routes.GET("/phone/:phone", httpTransport.viewByPhone)
	routes.PATCH("/:id", httpTransport.update)
	routes.DELETE("/:id", httpTransport.delete)
}

func (h HTTP) create(c *gin.Context) {
	//check that user is allowed to make this request

	var data createClientRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err); return
	}

	newClient := h.fromCreateClientRequest(&data)
	created, err := h.service.Create(c, newClient); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusCreated, clientCreated(created)); return
}

func (h HTTP) list(c *gin.Context) {
	clients, err := h.service.List(c); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetchedAll(clients)); return
}

func (h HTTP) viewById(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	fetchedClient, err := h.service.ViewById(c, id); if err != nil {
		PPA.Response(c, err); return
	}

	c.JSON(http.StatusOK, fetched(fetchedClient)); return
}

func (h HTTP) viewByPhone(c *gin.Context) {
	phone := c.Param("phone")
	if len(phone) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Phone Number Required")); return
	}

	fetchedClient, err := h.service.ViewByPhone(c, phone)if err != nil {
		PPA.Response(c, err); return
	}

	c.JSON(http.StatusOK, fetched(fetchedClient)); return
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
	c.JSON(http.StatusOK, clientUpdated(updated)); return
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

func deleted() gin.H {
	return gin.H{"success": true, "message": "User Deleted"}
}

func fetched(u *PPA.User) gin.H {
	return gin.H{"success": true, "message": "Fetched User", "payload": u}
}

func fetchedAll(u *[]PPA.User) gin.H {
	return gin.H{"success": true, "message": "Fetched User", "payload": u}
}

func userCreated(u *PPA.User) gin.H {
	return gin.H{"success": true, "message": "User Created", "payload": u}
}

func userUpdated(u *PPA.User) gin.H {
	return gin.H{"success": true, "message": "User Updated", "payload": u}
}
