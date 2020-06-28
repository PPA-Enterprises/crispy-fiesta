package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"github.com/PPA-Enterprises/crispy-fiesta/models"

)

type JobController struct {}

func (j *JobController) Create(c *gin.Context) {
	var data forms.SubmitJobCmd
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide relevant fields"})
		c.Abort()
		return
	}
	job := models.FromSubmitJobCmd(data)
	id, _ := job.PersistJob()
	c.JSON(201, gin.H{"message": id.String()})
}
