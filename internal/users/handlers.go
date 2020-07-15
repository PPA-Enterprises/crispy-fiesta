package users
import (
	"fmt"
	"context"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)
func Signup(c *gin.Context) {
	fmt.Println("got here")
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data signupUserCommand
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable,
		gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

	newUser, err := tryFromSignupUserCmd(&data); if err != nil {
		//500, hashing failed
		c.JSON(http.StatusInternalServerError,
			gin.H{"success": false, "message": err.Error()})
		return
	}
	id, err := newUser.signup(ctx); if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated,
		gin.H{"success": true, "payload": id.String(), "message": "User Created"});
}
