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

}
