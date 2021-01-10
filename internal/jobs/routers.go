package jobs

import(
	"github.com/gin-gonic/gin"
)

func JobRoutesRegister(router *gin.RouterGroup) {
	router.POST("/", submitJob)
	router.PATCH("/:id", update)
	router.GET("/:id", getJobByID)
	router.DELETE("/:id", deleteJobByID)
}
