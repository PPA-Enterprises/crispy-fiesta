package user

import(
	"PPA"
	"context"
	"pkg/common/mongo"
	dbQuery "pkg/api/user/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(*gin.Context, PPA.User) (*PPA.User, error)
	List(*gin.Context) (*[]PPA.User, error)
	ViewById(*gin.Context, string) (*PPA.User, error)
	//Delete(context.Context) error
	//Update(context.Context) (PPA.User, error)
}

func New(db *mongo.DBConnection, coll string, udb Repository, rbac RBAC, securer Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, securer: securer}
}

func Init(db *mongo.DBConnection, coll string, rbac RBAC, securer Securer) *User {
	return New(db, coll, dbQuery.User{}, rbac, securer)
}

type User struct {
	db *mongo.DBConnection
	collection string
	udb Repository
	rbac RBAC
	securer Securer
}

type Securer interface {
	Hash(string) string
}

type Repository interface {
	Create(*mongo.DBConnection, context.Context, *PPA.User) (*PPA.User, error)
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (*PPA.User, error)
	//Update(*mongo.DBConnection, PPA.User) (PPA.User, error)
	List(*mongo.DBConnection, context.Context) (*[]PPA.User, error)
	//Delete(*mongo.DBConnection) error
}

type RBAC interface {
	//User(*gin.Context) PPA.AuthUser
	//AccountCreate(*gin.Context) error
}
