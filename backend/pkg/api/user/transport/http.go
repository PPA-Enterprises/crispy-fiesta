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

func NewHTTP(service user.Service, router *gin.RouterGroup) {
	httpTransport := HTTP{service}
	routes := router.Group("/users")
	routes.POST("/", httpTransport.create)
	//routes.POST("/login", login)
	routes.GET("/", httpTransport.list)
	routes.GET("/email", httpTransport.viewByEmail)
	routes.GET("id/:id", httpTransport.viewById)
	//routes.PATCH("/:id", update)
	routes.DELETE("/:id",httpTransport.delete)
}

func (h HTTP) create(c *gin.Context) {
	//check that user is allowed to make this request

	var data signupRequest
	if c.BindJSON(&data) != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, bindFailure()); return
	}

	newUser := h.fromSignupRequest(&data)
	created, err := h.service.Create(c, newUser); if err != nil {
		//TODO: give proper error
		c.AbortWithStatusJSON(http.StatusNotAcceptable, bindFailure()); return
	}
	c.JSON(http.StatusCreated, userCreated(created)); return
}

func (h HTTP) viewById(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, bindFailure()); return
	}

	fetchedUser, err := h.service.ViewById(c, id); if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, fetchFailure()); return
	}

	c.JSON(http.StatusOK, fetched(fetchedUser)); return
}

func (h HTTP) list (c *gin.Context) {
	users, err := h.service.List(c); if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, fetchFailure()); return
	}
	c.JSON(http.StatusOK, fetchedAll(users)); return
}

func (h HTTP) viewByEmail(c *gin.Context) {
	var data emailRequest
	if c.BindJSON(&data) != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, bindFailure()); return
	}

	fetchedUser, err := h.service.ViewByEmail(c, h.fromEmailRequest(&data)); if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, fetchFailure()); return
	}

	c.JSON(http.StatusOK, fetched(fetchedUser)); return
}

func (h HTTP) delete(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, bindFailure()); return
	}

	if err := h.service.Delete(c, id); err != nil {
		//TODO: give proper error
		c.AbortWithStatusJSON(http.StatusNotAcceptable, bindFailure()); return
	}
	c.JSON(http.StatusOK, deleted()); return
}

func bindFailure() gin.H {
	return gin.H{"success": false, "message": "Provide relevant fields"}
}

func deleted() gin.H {
	return gin.H{"success": true, "message": "User Deleted"}
}

func fetchFailure() gin.H {
	return gin.H{"success": false, "message": "Failed to Fetch User"}
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
