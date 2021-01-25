package transport

import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pkg/api/user"
)

type signupRequest struct {
	Name string `json:"name" binding:"required,alphanum"`
	Email string `json:"email" binding:"required" email`
	Password string `json:"password" binding:"required,alphanum"`
}

type emailRequest struct {
	Email string `json:"email" binding:"required" email`
}

type updateRequest struct {
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty" email`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) user.Update {
	return user.Update {
		Name: data.Name,
		Email: data.Email,
	}
}

func (h HTTP) fromEmailRequest(data *emailRequest) string {
	return data.Email
}

func (h HTTP) fromSignupRequest(data *signupRequest) PPA.User {
	oid := primitive.NewObjectID()

	return PPA.User {
		ID: oid,
		Name: data.Name,
		Email: data.Email,
		Password: data.Password,
	}
}
