package users
import (
	"context"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)
func signup(c *gin.Context) {
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
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		return
	}

	id, err := newUser.signup(ctx); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated,
		gin.H{"success": true, "payload": id.String(), "message": "User Created"});
}

func login(c *gin.Context) {
	//ctx, cancel := context.WithTimeout(c, 5*time.Second)
	//defer cancel()

	var credentials loginUserCommand
	if c.BindJSON(&credentials) != nil {
		c.JSON(http.StatusNotAcceptable,
		gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

}
