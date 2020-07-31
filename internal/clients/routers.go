package clients

import (
	"github.com/gin-gonic/gin"
)

func ClientRoutesRegister(router *gin.RouterGroup) {
	router.GET("/phone/:phone", getClientByPhone)
	router.PUT("/", update)
	router.GET("/search", fuzzyClientSearch)
	router.GET("/:id", getClientByID)
}
