package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"github.com/PPA-Enterprises/crispy-fiesta/models"
)

type ClientController struct{}

var clientModel = new(models.ClientModel)

func (createClient *ClientController) CreateClient(c *gin.Context) {
	var data forms.CreateClientCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide all required fields"})
		c.Abort()
		return
	}

	result, err := clientModel.CreateClient(data)

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem creating client."})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Created client successfully", "_id": result.InsertedID})

}

func (getClient *ClientController) GetClientById(c *gin.Context) {
	id := c.Param("id")
	var client models.Client

	client, err := clientModel.GetClientById(id)

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem getting client", "error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Got client successfully", "client": client})
}
