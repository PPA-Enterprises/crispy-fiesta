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
		c.Abort()
		return
	}

	id, err := newUser.signup(ctx); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated,
		gin.H{"success": true, "payload": id.String(), "message": "User Created"});
}

func login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var credentials loginUserCommand
	if c.BindJSON(&credentials) != nil {
		c.JSON(http.StatusNotAcceptable,
		gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

	jwt, err := authenticate(ctx, credentials); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK,
		gin.H{"success": true, "payload": jwt, "message": "User Authenticated"});

}

func getUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	users, err := fetchUsers(ctx)
	if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}
	if len(users) > 0 {
		c.JSON(http.StatusOK,
			gin.H{"success": true, "payload": users})
		return
	}
		empty := make([]string, 0)
		c.JSON(http.StatusOK,
			gin.H{"success": true, "payload": empty})
}

func update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data userUpdateCommand
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable,
		gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(http.StatusBadRequest,
		gin.H{"success": false, "message": "Provide an id"})
		c.Abort()
		return
	}

	updatePatch, err := tryFromUpdateUserCmd(&data, id); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	updatedDoc, patchErr := updatePatch.patch(ctx, false); if patchErr != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": patchErr.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated,
	gin.H{"success": true, "payload": updatedDoc, "message": "User Updated"});
}

func delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(http.StatusBadRequest,
		gin.H{"success": false, "message": "Provide an id"})
		c.Abort()
		return
	}

	if err := deleteUser(ctx, id); err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK,
	gin.H{"success": true, "message": "User Removed"});
}
