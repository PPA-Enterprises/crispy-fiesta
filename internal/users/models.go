package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	passwordUtils "internal/common"
)

type userModel struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	IsVerified bool `json:"is_verified" bson:"is_verified"`
}

func tryFromSignupUserCmd(data signupUserCommand) (*userModel, error) {
	encrypted, err := passwordUtils.HashFromPlaintext(data.Password)
	if err != nil {
		return nil, err
	}

	return &userModel{
		Name: data.Name,
		Email: data.Email,
		Password: encrypted,
	}, nil
}

/*func (u *userModel) signup() (*ID, error) {
	//coll := dbConnect.Use(name, "user")
	// if user exists, return err
	// insert user, return id
}*/
