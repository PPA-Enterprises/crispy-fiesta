package users

type signupUserCommand struct {
	Name		string `json:"name" binding:"required"`
	Email		string `json:"email" binding:"required"`
	Password	string `json:"password" binding:"required"`
}

type loginUserCommand struct {
	Email		string `json:"email" binding:"required,email"`
	Password	string `json:"password" binding:"required"`
}

type userUpdateCommand struct {
	ID			string `json:"_id" binding:"required"`
	Name		string `json:"name" binding:"required"`
	Email		string `json:"email" binding:"required"`
}
