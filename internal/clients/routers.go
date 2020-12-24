package clients

import (
	"github.com/gin-gonic/gin"
)

func ClientRoutesRegister(router *gin.RouterGroup) {
	router.GET("/client/:id", getClientByID)
	router.GET("/phone/:phone", getClientByPhone)
	router.PUT("/", update)
	router.GET("/", getLatestClients)
	router.GET("/search/:query", fuzzyClientSearch)
}
