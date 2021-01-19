package user

import(
	"PPA"
	"context"
	"pkg/common/mongo"
	dbQuery "pkg/api/user/infra/mongo"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(context.Context, PPA.User) (*PPA.User, error)
	//List(context.Context) ([]PPA.User, error)
	//View(context.Context) (PPA.User, error)
	//Delete(context.Context) error
	//Update(context.Context) (PPA.User, error)
}

func New(db *mongo.DBConnection, udb Repository, rbac RBAC, securer Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, securer: securer}
}

func Init(db *mongo.DBConnection, rbac RBAC, securer Securer) *User {
	return New(db, dbQuery.User{}, rbac, securer)
}

type User struct {
	db *mongo.DBConnection
	udb Repository
	rbac RBAC
	securer Securer
}

type Securer interface {
	Hash(string) string
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.User) (*PPA.User, error)
	//View(*mongo.DBConnection, PPA.User) (PPA.User, error)
	//Update(*mongo.DBConnection, PPA.User) (PPA.User, error)
	//List(*mongo.DBConnection) ([]PPA.User, error)
	//Delete(*mongo.DBConnection) error
}

type RBAC interface {
	User(*gin.Context) PPA.AuthUser
	AccountCreate(*gin.Context) error
}
