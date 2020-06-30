package controllers

import (
	"net/http"

	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"github.com/PPA-Enterprises/crispy-fiesta/models"
	"github.com/gin-gonic/gin"
)

var jobModel = new(models.JobModel)

type JobController struct{}

func (j *JobController) Create(c *gin.Context) {
	var data forms.SubmitJobCmd
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable,
			gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}
	job := models.FromSubmitJobCmd(data)
	id, err := job.PersistJob()

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated,
		gin.H{"success": true, "payload": id.String(), "message": "Job Created"})
}

func (j *JobController) Update(c *gin.Context) {
	var data forms.UpdateJobCmd

	if c.BindJSON(&data) != nil {
		c.JSON(406,
			gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

	_, err := jobModel.UpdateJob(data)

	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(202, gin.H{"success": true, "message": "Account updated."})

}
