package api

import (
	"pkg/common/config"
	"pkg/common/mongo"
	"pkg/common/rbac"
	"pkg/common/secure"
	userTransport "pkg/api/user/transport"
	userService "pkg/api/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Configuration) error {
	db := mongo.Init("mongodb://localhost:27017")
	server := gin.Default()
	server.Use(cors.Default())
	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found."})
	})

	v1 := server.Group("/api/v1")
	rbac := rbac.Service{}
	security := secure.New()
	userTransport.NewHTTP(userService.Init(db, rbac, security), v1)
	server.Run()
	return nil
}
