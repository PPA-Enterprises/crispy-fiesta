package jobs

import (
	"context"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)

func submitJob(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data submitJobCmd
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable,
		gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

	newJob := fromSubmitJobCmd(&data)
	oid, err := newJob.create(ctx); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated,
		gin.H{"success": true, "payload": oid.String(), "message": "Job Created"})
}

func update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data updateJobCmd
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable,
		gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

	update, err := tryFromUpdateJobCmd(&data); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	err = update.put(ctx, false); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusAccepted,
		gin.H{"success": true, "message": "Job Updated"})
}

func getJobByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	id := c.Param("id") //returns empty string if not there
	if len(id) <= 0 {
		c.JSON(http.StatusBadRequest,
		gin.H{"success": false, "message": "Provide an id"})
		c.Abort()
		return
	}

	job, err := jobByID(ctx, id); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"success": true, "payload": job})
}

