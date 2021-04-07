package api

import (
	"PPA"
	"crypto/sha1"
	"pkg/common/config"
	"pkg/common/mongo"
	"pkg/common/rbac"
	"pkg/common/secure"
	"pkg/common/jwt"
	"pkg/common/eventlog"
	authMw "pkg/common/middleware/auth"
	userTransport "pkg/api/user/transport"
	userService "pkg/api/user"
	userRepo "pkg/api/user/infra/mongo"
	authTransport "pkg/api/auth/transport"
	authService "pkg/api/auth"
	clientTransport "pkg/api/client/transport"
	clientService "pkg/api/client"
	clientRepo "pkg/api/client/infra/mongo"
	jobTransport "pkg/api/job/transport"
	jobService "pkg/api/job"
	jobRepo "pkg/api/job/infra/mongo"
	eventlogTransport "pkg/api/eventlog/transport"
	eventlogService "pkg/api/eventlog"
	tinterTransport "pkg/api/tinter/transport"
	tinterService "pkg/api/tinter"
	tinterRepo "pkg/api/tinter/infra/mongo"
	clientLabelTransport "pkg/api/clientlabel/transport"
	clientLabelService "pkg/api/clientlabel"
	labelRepo "pkg/api/clientlabel/infra/mongo"
	//jobRepo "pkg/api/job/infra/mongo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Configuration) error {
	db := mongo.Init("mongodb://localhost:27017/?replicaSet=replicaset&readPreference=primaryPreferred", "PPA")
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

func StreamAPI(cfg *config.Configuration) error {
	db := mongo.Init("mongodb://localhost:27017/?replicaSet=replicaset&readPreference=primaryPreferred", "PPA")
	server := gin.Default()
	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found."})
	})

	jobStream := NewStreamServer()
	clientStream := NewStreamServer()
	server.Use(jobStream.ServeHTTP())
	server.Use(clientStream.ServeHTTP())

	v1 := server.Group("/stream/v1")
	rbac := rbac.Service{}
	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, "my asdfasfasdfsafsadfasfasfasfsafasfsadfsadsecret from os get env", cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)

	if err != nil {
		return err
	}

	authMiddleware := authMw.Middleware(jwt)
	jobTransport.NewStream(jobService.InitStream(db, rbac, tinterRepo.Tinter{}), jobStream, v1, authMiddleware)
	clientTransport.NewStream(clientService.InitStream(db, rbac, jobRepo.Job{}, labelRepo.ClientLabel{}), clientStream, v1, authMiddleware)
	server.Run(":8085")
	return nil
}


func NewStreamServer() (event *PPA.StreamEvent) {
	event = &PPA.StreamEvent{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.Listen()

	return
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
