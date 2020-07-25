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

