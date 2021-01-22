package auth

import (
	"PPA"
	"github.com/gin-gonic/gin"
	"context"
	"pkg/common/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	dbQuery "pkg/api/auth/infra/mongo"
)

type Auth struct {
	db *mongo.DBConnection
	udb UDB
	tokenGen TokenGenerator
	securer Securer
	rbac RBAC
}

func New(db *mongo.DBConnection, udb UDB, tkgen TokenGenerator, sec Securer, rbac RBAC) Auth {
	return Auth {
		db: db,
		udb: udb,
		tokenGen: tkgen,
		securer: sec,
		rbac: rbac,
	}
}

func Init(db *mongo.DBConnection, tkgen TokenGenerator, sec Securer, rbac RBAC) Auth {
	return New(db, dbQuery.User{}, tkgen, sec, rbac)
}

type Service interface {
	Authenticate(*gin.Context, string, string) (PPA.AuthToken, error)
	Refresh(*gin.Context, string) (string, error)
}

type UDB interface {
	ViewById(*mongo.DBConnection, context.Context, primitive.ObjectID) (PPA.user, error)
	ViewByEmail(*mongo.DBConnection, context.Context, string) (*PPA.User, error)
	ViewByToken(*mongo.DBConnection, context.Context, string) (*PPA.User, error)
	Update(*mongo.DBConnection, context.Context, *PPA.User) error

}

type TokenGenerator interface {
	GenerateToken(PPA.User) (string, error)
}

type Securer interface {
	HashMatchesPassword(string, string) bool
	Token(string) string
}

type RBAC interface {}
