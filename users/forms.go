package users

type signupUserCommand struct {
	Name		string `json:"name" binding:"required"`
	Email		string `json:"name" binding:"required"`
	Password	string `json:"name" binding:"required"`
}
