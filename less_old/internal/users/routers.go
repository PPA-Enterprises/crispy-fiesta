package users

import (
	"github.com/gin-gonic/gin"
)

func UserRoutesRegister(router *gin.RouterGroup) {
	router.POST("/", signup)
	router.POST("/login", login)
	router.GET("/", getUsers)
	router.PATCH("/:id", update)
	router.DELETE("/:id", delete)
}
