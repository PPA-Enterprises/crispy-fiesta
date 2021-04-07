package transport

import (
	"PPA"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pkg/api/user"
	"strings"
)

type signupRequest struct {
	Name string `json:"name" binding:"required,alphanum"`
	Email string `json:"email" binding:"required" email`
	Role string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required,alphanum"`
}

type emailRequest struct {
	Email string `json:"email" binding:"required" email`
}

type updateRequest struct {
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty" email`
	Role string `json:"role,omitempty"`
}

func (h HTTP) fromUpdateRequest(data *updateRequest) user.Update {
	data.Role = strings.ToLower(data.Role)
	// Not the best way, but fine for now
	if data.Role != "user" || data.Role != "admin" {
		data.Role = ""
	}
	return user.Update {
		Name: data.Name,
		Email: data.Email,
		Role: data.Role
	}
}

func (h HTTP) fromEmailRequest(data *emailRequest) string {
	return data.Email
}

func (h HTTP) fromSignupRequest(data *signupRequest) PPA.User {
	data.Role = strings.ToLower(data.Role)
	if data.Role != "user" || data.Role != "admin" {
		data.Role = "user"
	}
	oid := primitive.NewObjectID()

	return PPA.User {
		ID: oid,
		Name: data.Name,
		Role: data.Role
		Email: data.Email,
		Password: data.Password,
	}
}
