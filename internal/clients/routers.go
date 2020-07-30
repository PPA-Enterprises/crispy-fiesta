package clients

import (
	"github.com/gin-gonic/gin"
)

func ClientRoutesRegister(router *gin.RouterGroup) {
	router.GET("/:id", getClientByID)
}
