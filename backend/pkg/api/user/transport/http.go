package transport

import(
	"PPA"
	"net/http"
	"pkg/api/user"
	"github.com/gin-gonic/gin"
)

type HTTP struct {
	service user.Service
}

const (
	BadRequest = http.StatusBadRequest
)

func NewHTTP(service user.Service, router *gin.RouterGroup, mw gin.HandlerFunc) {
	httpTransport := HTTP{service}
	routes := router.Group("/users")
	routes.POST("/", httpTransport.create)
	routes.GET("/", mw, httpTransport.list)
	routes.GET("/email", httpTransport.viewByEmail)
	routes.GET("id/:id", httpTransport.viewById)
	routes.PATCH("/:id", httpTransport.update)
	routes.DELETE("/:id",httpTransport.delete)
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
	c.JSON(http.StatusCreated, userCreated(created)); return
}

func (h HTTP) viewById(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "ID Required")); return
	}

	fetchedUser, err := h.service.ViewById(c, id); if err != nil {
		PPA.Response(c, err); return
	}

	c.JSON(http.StatusOK, fetched(fetchedUser)); return
}

func (h HTTP) list (c *gin.Context) {
	users, err := h.service.List(c); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, fetchedAll(users)); return
}

func (h HTTP) viewByEmail(c *gin.Context) {
	var data emailRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		PPA.Response(c, err); return
	}

	fetchedUser, err := h.service.ViewByEmail(c, h.fromEmailRequest(&data)); if err != nil {
		PPA.Response(c, err); return
	}

	c.JSON(http.StatusOK, fetched(fetchedUser)); return
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
	c.JSON(http.StatusOK, userUpdated(updated)); return
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
