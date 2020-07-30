package main

import (
	//"github.com/PPA-Enterprises/crispy-fiesta/controllers"
	//"github.com/PPA-Enterprises/crispy-fiesta/helpers"
	"internal/db"
	userRoutes "internal/users"
	jobRoutes "internal/jobs"
	clientRoutes "internal/clients"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	router := gin.Default()
	conn := db.Init("mongodb://localhost:27017")
	defer conn.Disconnect()

	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	userRoutes.UserRoutesRegister(v1.Group("/users"))
	jobRoutes.JobRoutesRegister(v1.Group("/jobs"))
	clientRoutes.ClientRoutesRegister(v1.Group("/clients"))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found."})
	})

	router.Run()
}

