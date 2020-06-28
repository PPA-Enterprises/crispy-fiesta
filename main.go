package main

import (
	"github.com/PPA-Enterprises/crispy-fiesta/controllers"
	"github.com/PPA-Enterprises/crispy-fiesta/helpers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		hello := new(controllers.HelloWorldController)
		user := new(controllers.UserController)
		job := new(controllers.JobController)

		v1.GET("/hello", TokenAuthMiddleware(), hello.Default)
		v1.POST("/signup", user.Signup)
		v1.POST("/login", user.Login)
		v1.POST("/job", job.Create)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found."})
	})

	router.Run()
}

//middlewares
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helpers.TokenValid(c.Request)
		if err != nil {
			c.JSON(401, gin.H{"message": "Not authenticated."})
			c.Abort()
			return
		}
		c.Next()
	}
}
