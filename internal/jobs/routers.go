package jobs

import(
	"github.com/gin-gonic/gin"
)

func JobRoutesRegister(router *gin.RouterGroup) {
	router.POST("/", submitJob)
	router.PUT("/", update)
}
