package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/PPA-Enterprises/crispy-fiesta/models"
)

type ClientController struct{}

var clientModel = new(models.ClientModel)

func (getClient *ClientController) GetClientById(c *gin.Context) {
	id := c.Param("id")
	var client models.Client

	client, err := clientModel.GetClientById(id)

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem getting client"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Got client successfully", "client": client})
}
