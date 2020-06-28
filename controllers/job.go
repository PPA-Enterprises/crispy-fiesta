package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"github.com/PPA-Enterprises/crispy-fiesta/models"

)

type JobController struct {}

func (j *JobController) Create(c *gin.Context) {
	var data forms.SubmitJobCmd
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Provide relevant fields"});
		c.Abort();
		return
	}
	job := models.FromSubmitJobCmd(data);
	id, err := job.PersistJob();

	if err != nil {
		//TODO: check for err type and respond accordingly
		c.JSON(http.StatusCreated,
			gin.H{"success": false, "payload": "", "message": "Job Create Failed"});
		return;
	}
	c.JSON(http.StatusCreated,
		gin.H{"success": true, "payload": id.String(), "message": "Job Created"});
}
