package transport

import (
	"context"
	"PPA"
	"pkg/api/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type signupRequest struct {
	Name string `json:"name" binding:"required,alphanum"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum"`
}

func (h HTTP) fromSignupRequest(ctx context.Context, data *signupRequest) (PPA.User, error) {
	encrypted := h.service.securer(data.Password)
	oid := primitive.NewObjectID()

	for oidExists(ctx, oid) {
		oid = primitive.NewObjectID()
	}
	return PPA.User {
		ID: oid,
		Name: data.Name,
		Email: data.Email,
		Password: encrypted,
		IsDeleted: false,
	}
}

func (h HTTP) oidExists(ctx context.Context, oid primitive.ObjectID) bool {
	coll := h.service.db.Use(db.DefaultDatabase, "users")

	var user PPA.User
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&user)
	return err == nil
}
