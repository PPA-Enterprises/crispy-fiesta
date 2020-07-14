package users

type user struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	IsVerified bool `json:"is_verified" bson:"is_verified"`
}

// TODO: hash password here
func tryFromSubmitJobCmd(data signupUserCommand) (*user, error) {
	return &user{
		Name: data.Name,
		Email: data.Email,
		Password: data.Password,
	}, nil
}

func (u *user) signup() (*ID, error) {
	//coll := dbConnect.Use(name, "user")
	// if user exists, return err
	// insert user, return id
}
