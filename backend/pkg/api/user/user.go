package user

import (
	"context"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"PPA"
)

const (
	NotFound = http.StatusNotFound
)

var OidNotFound = PPA.NewAppError(NotFound, "Does not exist")
const (
	Collection = "users"
	EventTag = "m"
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
	created, err := u.udb.Create(u.db, ctx, &req); if err != nil {
		return nil, err
	}

	editor := PPA.Editor {
		OID: created.ID,
		Name: "Bob",
		Collection: "Bob" + created.ID.Hex() + "a",

	}
	created.AppendLog(u.eventLogger.LogCreated(ctx, u.eventLogger.GenerateEvent(created, EventTag), editor))
	u.udb.LogEvent(u.db, ctx, created)

	return created, nil
}

func (u User) ViewById(c *gin.Context, id string) (*PPA.User, error) {
	//additional security stuff
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, OidNotFound
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

func (u User) Delete(c *gin.Context, id string) error {
	//additional security stuff?
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return OidNotFound
	}

	return u.udb.Delete(u.db, ctx, oid)
}

type Update struct {
	Name string
	Email string
}

func (u User) Update(c *gin.Context, req Update, id string) (*PPA.User, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, OidNotFound
	}

	if err := u.udb.Update(u.db, ctx, oid, &PPA.User {
		ID: primitive.NilObjectID,
		Name: req.Name,
		Email: req.Email,
		Password: "",
	}); err != nil {
		return nil, err
	}
	return u.udb.ViewById(u.db, ctx, oid)
}

func (u User) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := u.db.Use(Collection)

	var inserted PPA.User
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&inserted)
	return err == nil
}
