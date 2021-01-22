package transport

import (
	"net/http"
	"pkg/api/auth"
	"github.com/gin-gonic/gin"
)

type HTTP struct {
	service auth.Service
}

func NewHTTP(service auth.Service, router *gin.RouterGroup) {
	httpTransport := HTTP{service}
	routes := router.Group("/auth")
	routes.POST("/", httpTransport.login)
	routes.Get("/refresh/:token", httpTransport.refresh)
}

func (h HTTP) login(c *gin.Context) {
	var cred credentials
	if err := c.ShouldBindJSON(&cred); err != nil {
		PPA.Response(c, err); return
	}

	authToken, err := h.service.Authenticate(c, cred.Email, cred.Password); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, authenticated(authToken)); return
}

func (h HTTP) refresh(c *gin.Context) error {
	token := c.Param("token")
	if len(token) <= 0 {
		PPA.Response(c, PPA.NewAppError(BadRequest, "Token Required")); return
	}

	token, err := h.service.Refresh(c, token); if err != nil {
		PPA.Response(c, err); return
	}
	c.JSON(http.StatusOK, refreshed(token)); return
}

func authenticated(t *PPA.AuthToken) gin.H {
	return gin.H{"success": true, "message": "User Logged In", "payload": t}
}

func refreshed(t string) gin.H {
	return gin.H{"success": true, "message": "", "payload": t}
}
