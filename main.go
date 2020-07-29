package main

import (
	//"github.com/PPA-Enterprises/crispy-fiesta/controllers"
	//"github.com/PPA-Enterprises/crispy-fiesta/helpers"
	"internal/db"
	userRoutes "internal/users"
	jobRoutes "internal/jobs"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//conn := db.Init("mongodb://localhost:27017?replicaSet=myRepl")
	conn := db.Init("mongodb://localhost:27017", "myrepl")
	defer conn.Disconnect()

	/*v1 := router.Group("/api/v1")
	{
		hello := new(controllers.HelloWorldController)
		user := new(controllers.UserController)
		job := new(controllers.JobController)
		client := new(controllers.ClientController)

		v1.GET("/hello", TokenAuthMiddleware(), hello.Default)
		v1.POST("/signup", user.Signup)
		v1.POST("/login", user.Login)
		v1.POST("/job", job.CreateJob)
		v1.PUT("/job", job.Update)
		v1.POST("/client", client.CreateClient)
		v1.GET("/client/:id", client.GetClientById)
		v1.GET("/client", client.GetAllClients)
	}*/

	v1 := router.Group("/api/v1")
	userRoutes.UserRoutesRegister(v1.Group("/users"))
	jobRoutes.JobRoutesRegister(v1.Group("/jobs"))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found."})
	})

	router.Run()
}
/*
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
}*/
