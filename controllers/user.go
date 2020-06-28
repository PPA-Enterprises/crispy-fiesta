package controllers

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"github.com/PPA-Enterprises/crispy-fiesta/helpers"
	"github.com/PPA-Enterprises/crispy-fiesta/models"
)

var userModel = new(models.UserModel)

type UserController struct{}

func (u *UserController) Signup(c *gin.Context) {
	var data forms.SignupUserCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide relevant fields"})
		c.Abort()
		return
	}

	_, err := userModel.Signup(data)

	log.Print(err)

	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "New user account registered"})
}

func (u *UserController) Login(c *gin.Context) {
	var data forms.LoginUserCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide relevant fields"})
		c.Abort()
		return
	}

	foundUser, err := userModel.Login(data)

	if err != nil {
		c.JSON(401, gin.H{"message": "User Not Found."})
		c.Abort()
		return
	}

	match, err := helpers.ComparePasswordAndHash(data.Password, foundUser.Password)

	if err != nil {
		c.JSON(400, gin.H{"message": "Error with password compare."})
		c.Abort()
		return
	}

	if match {
		//generate token
		//status 200
		c.JSON(200, gin.H{"message": "Logged in!!"})
	} else {
		//wrong password
		//status 401
		c.JSON(401, gin.H{"message": "Password incorrect."})
		c.Abort()
		return
	}

}
