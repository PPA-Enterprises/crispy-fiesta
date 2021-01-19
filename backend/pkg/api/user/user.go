package user

import (
	"context"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"PPA"
)

func (u User) Create(c *gin.Context, req PPA.User) (*PPA.User, error) {
	//additional security stuff like if user is allowed to do this
	/*if err := self.rbac.AccountCreate(c); err != nil {
		return PPA.User{}, err
	}*/
	// Potentially extract deadline from *gin.Context.Request.Context
	// https://golang.org/pkg/context/#WithDeadline
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	// ensure OID is unique
	for u.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}

	req.Password = u.securer.Hash(req.Password)
	return u.udb.Create(u.db, ctx, &req)
}

func (u User) List() {}
func (u User) View() {}
func (u User) Delete() {}

func (u User) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := u.db.Use("PPA", "users")

	var inserted PPA.User
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}
