package users

import (
	"github.com/gin-gonic/gin"
)

func UserRoutesRegister(router *gin.RouterGroup) {
	router.POST("/", Signup)
}
