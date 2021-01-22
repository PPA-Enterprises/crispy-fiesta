package transport

type credentials struct {
	Email string `json:"email" binding:"required" email`
	Password string `json:"password" binding:"required"`
}
