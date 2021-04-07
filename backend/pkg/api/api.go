package api

import (
	"crypto/sha1"
	authService "pkg/api/auth"
	authTransport "pkg/api/auth/transport"
	clientService "pkg/api/client"
	clientRepo "pkg/api/client/infra/mongo"
	clientTransport "pkg/api/client/transport"
	clientLabelService "pkg/api/clientlabel"
	labelRepo "pkg/api/clientlabel/infra/mongo"
	clientLabelTransport "pkg/api/clientlabel/transport"
	eventlogService "pkg/api/eventlog"
	eventlogTransport "pkg/api/eventlog/transport"
	jobService "pkg/api/job"
	jobRepo "pkg/api/job/infra/mongo"
	jobTransport "pkg/api/job/transport"
	tinterService "pkg/api/tinter"
	tinterRepo "pkg/api/tinter/infra/mongo"
	tinterTransport "pkg/api/tinter/transport"
	userService "pkg/api/user"
	userRepo "pkg/api/user/infra/mongo"
	userTransport "pkg/api/user/transport"
	"pkg/common/config"
	"pkg/common/eventlog"
	"pkg/common/jwt"
	authMw "pkg/common/middleware/auth"
	"pkg/common/mongo"
	"pkg/common/rbac"
	"pkg/common/secure"
	"time"

	//jobRepo "pkg/api/job/infra/mongo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Configuration) error {
	db := mongo.Init("mongodb://localhost:27017", "PPA")
	server := gin.Default()
	//server.Use(cors.Default())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://ppaenterprises.com"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "OPTIONS", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://ppaenterprises.com"
		},
		MaxAge: 12 * time.Hour,
	}))

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
	logger := eventlog.New(db)

	authTransport.NewHTTP(authService.Init(db, jwt, security, rbac, userRepo.User{}), v1)
	userTransport.NewHTTP(userService.Init(db, rbac, security, logger), v1, authMiddleware)
	clientTransport.NewHTTP(clientService.Init(db, rbac, jobRepo.Job{}, labelRepo.ClientLabel{}, logger), v1, authMiddleware)
	jobTransport.NewHTTP(jobService.Init(db, rbac, clientRepo.Client{}, tinterRepo.Tinter{}, logger), v1, authMiddleware)
	tinterTransport.NewHTTP(tinterService.Init(db, rbac, jobRepo.Job{}, logger), v1, authMiddleware)
	clientLabelTransport.NewHTTP(clientLabelService.Init(db, rbac, logger), v1, authMiddleware)
	eventlogTransport.NewHTTP(eventlogService.Init(db, userRepo.User{}, rbac), v1, authMiddleware)
	server.Run(cfg.Server.Port)
	return nil
}
