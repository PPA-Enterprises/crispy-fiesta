package transport

import(
	"context"
	"time"
	"net/http"
	"strconv"
	"pkg/api/user"
	"github.com/gin-gonic/gin"
)

type HTTP struct {
	service user.Service
}

func NewHTTP(service user.Service, router *gin.RouterGroup) {
	httpService := HTTP{service}
	routes := router.Group("/users")
	routes.POST("/", httpService.create)
	//routes.POST("/login", login)
	//routes.GET("/", getUsers)
	//routes.PATCH("/:id", update)
	//routes.DELETE("/:id", delete)
}

func (h HTTP) create(c *gin.Context) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data signupRequest
	if c.BindJSON(&data) != nil {
		return c.AbortWithStatusJSON(http.StatusNotAcceptable, bindFailure())
	}

	newUser := h.fromSignupRequest(ctx, data)
	if user, err := h.service.Create(ctx, h.db, &newUser); err != nil {
		//return error
	}
	return c.JSON(http.StatusCreated, userCreated(user))
}

func bindFailure() gin.H {
	return gin.H{"success": false, "message": "Provide relevant fields"}
}

func userCreated(user *PPA.User) gin.H {
	return gin.H{"success": true, "message": "User Created", "payload": user}
}
