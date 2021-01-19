package transport

import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type signupRequest struct {
	Name string `json:"name" binding:"required,alphanum"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum"`
}

func (h HTTP) fromSignupRequest(data *signupRequest) PPA.User {
	oid := primitive.NewObjectID()

	return PPA.User {
		ID: oid,
		Name: data.Name,
		Email: data.Email,
		Password: data.Password,
		IsDeleted: false,
	}
}
