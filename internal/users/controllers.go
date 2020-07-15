package users
import (
	"fmt"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)
func Signup(c *gin.Context) {
	fmt.Println("got here")
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data signupUserCommand
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"msg": ""})
		c.Abort()
		return
	}

	newUser, err := tryFromSignupUserCmd(&data); if err != nil {
		//500, hashing failed
	}
	newUser.signup(ctx)
	c.JSON(201, gin.H{"msg":""})
}
