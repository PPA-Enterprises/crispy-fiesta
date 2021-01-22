package user

import (
	"context"
	"errors"
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
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	// ensure OID is unique
	for u.oidExists(ctx, req.ID) {
		req.ID = primitive.NewObjectID()
	}

	req.Password = u.securer.Hash(req.Password)
	return u.udb.Create(u.db, ctx, &req)
}

func (u User) ViewById(c *gin.Context, id string) (*PPA.User, error) {
	//additional security stuff
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, errors.New("")
	}

	return u.udb.ViewById(u.db, ctx, oid)
}

func (u User) List(c *gin.Context) (*[]PPA.User, error) {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	//Populate anything?
	return u.udb.List(u.db, ctx)
}

func (u User) ViewByEmail(c *gin.Context, email string) (*PPA.User, error) {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	return u.udb.ViewByEmail(u.db, ctx, email)
}

func (u User) Delete() {}

func (u User) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := u.db.Use(u.collection)

	var inserted PPA.User
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}
