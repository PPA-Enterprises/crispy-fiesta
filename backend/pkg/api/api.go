package api

import (
	"crypto/sha1"
	"pkg/common/config"
	"pkg/common/mongo"
	"pkg/common/rbac"
	"pkg/common/secure"
	"pkg/common/jwt"
	authMw "pkg/common/middleware/auth"
	userTransport "pkg/api/user/transport"
	userService "pkg/api/user"
	userRepo "pkg/api/user/infra/mongo"
	authTransport "pkg/api/auth/transport"
	authService "pkg/api/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Configuration) error {
	db := mongo.Init("mongodb://localhost:27017", "PPA")
	server := gin.Default()
	server.Use(cors.Default())
	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found."})
	})

	v1 := server.Group("/api/v1")
	rbac := rbac.Service{}
	security := secure.New(sha1.New())
	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, "my asdfasfasdfsafsadfasfasfasfsafasfsadfsadsecret from os get env", cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)

	if err != nil {
		return err
	}

	authMiddleware := authMw.Middleware(jwt)

	userTransport.NewHTTP(userService.Init(db, "users", rbac, security), v1, authMiddleware)
	authTransport.NewHTTP(authService.Init(db, jwt, security, rbac, userRepo.User{}), v1)
	server.Run()
	return nil
}
