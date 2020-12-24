package controllers

import (
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"github.com/PPA-Enterprises/crispy-fiesta/models"
	"github.com/gin-gonic/gin"
)

var jobModel = new(models.JobModel)

type JobController struct{}

func (createJob *JobController) CreateJob(c *gin.Context) {
	var data forms.SubmitJobCmd

	if c.BindJSON(&data) != nil {
		c.JSON(406,
			gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

	result, err := jobModel.CreateJob(data)

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem creating client."})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Created Job successfully", "_id job": result[0].InsertedID, "_id client": result[1].InsertedID})

}

// func (j *JobController) Create(c *gin.Context) {
// 	var data forms.SubmitJobCmd
// 	if c.BindJSON(&data) != nil {
// 		c.JSON(http.StatusNotAcceptable,
// 			gin.H{"success": false, "message": "Provide relevant fields"})
// 		c.Abort()
// 		return
// 	}
// 	job := models.FromSubmitJobCmd(data)
// 	id, err := job.PersistJob()

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError,
// 			gin.H{"success": false, "message": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated,
// 		gin.H{"success": true, "payload": id.String(), "message": "Job Created"})
// }

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
