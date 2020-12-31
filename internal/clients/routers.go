package clients

import (
	"github.com/gin-gonic/gin"
)

func ClientRoutesRegister(router *gin.RouterGroup) {
	router.GET("/client/:id", getClientByID)
	router.GET("/phone/:phone", getClientByPhone)
	router.PUT("/", update)
	router.GET("/", getClients)
	router.GET("/search/", fuzzyClientSearch)
	router.POST("/", createClient)
}
