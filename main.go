package main

import (
	"github.com/PPA-Enterprises/crispy-fiesta/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		hello := new(controllers.HelloWorldController)
		user := new(controllers.UserController)

		v1.GET("/hello", hello.Default)
		v1.POST("/signup", user.Signup)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found."})
	})

	router.Run()
}
