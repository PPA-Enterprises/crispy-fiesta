package api

import (
	//"crypto/sha1"
	//"os"

	//"github.com/ribice/gorsk/pkg/utl/zlog"
	//"github.com/ribice/gorsk/pkg/api/auth"
	//al "github.com/ribice/gorsk/pkg/api/auth/logging"
	//at "github.com/ribice/gorsk/pkg/api/auth/transport"
	//"github.com/ribice/gorsk/pkg/api/password"
	//pl "github.com/ribice/gorsk/pkg/api/password/logging"
	//pt "github.com/ribice/gorsk/pkg/api/password/transport"
	//"github.com/ribice/gorsk/pkg/api/user"
//	ul "github.com/ribice/gorsk/pkg/api/user/logging"
//	ut "github.com/ribice/gorsk/pkg/api/user/transport"

	//"github.com/ribice/gorsk/pkg/utl/config"
	"pkg/common/config"
	"pkg/common/mongo"
//	"github.com/ribice/gorsk/pkg/utl/jwt"
//	authMw "github.com/ribice/gorsk/pkg/utl/middleware/auth"
//	"github.com/ribice/gorsk/pkg/utl/rbac"
//	"github.com/ribice/gorsk/pkg/utl/secure"
//	"github.com/ribice/gorsk/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	_ = mongo.Init("mongodb://localhost:27017")
/*
	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rbac := rbac.Service{}
	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	authMiddleware := authMw.Middleware(jwt)

	at.NewHTTP(al.New(auth.Initialize(db, jwt, sec, rbac), log), e, authMiddleware)

	v1 := e.Group("/v1")
	v1.Use(authMiddleware)

	ut.NewHTTP(ul.New(user.Initialize(db, rbac, sec), log), v1)
	pt.NewHTTP(pl.New(password.Initialize(db, rbac, sec), log), v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})
*/
	return nil
}
