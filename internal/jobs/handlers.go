package jobs

import (
	"context"
	"time"
	"net/http"
	eventLogTypes "internal/event_log/types"
	"github.com/gin-gonic/gin"
)

func submitJob(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data submitJobCmd
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable,
		gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort(); return
	}

	newJob := fromSubmitJobCmd(&data)
	editor := eventLogTypes.Editor {
		Oid: newJob.ID,
		Name: "Bob",
		Collection: "Bob123",
	}
	oid, err := newJob.create(ctx, &editor); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
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
		c.Abort(); return
	}

	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(http.StatusBadRequest,
		gin.H{"success": false, "message":"Provide and id"})
		c.Abort(); return
	}

	update, err := tryFromUpdateJobCmd(&data, id); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
	}

	editor := eventLogTypes.Editor {
		Oid: update.ID,
		Name: "Bob",
		Collection: "Bob123",
	}
	updated, err := update.Patch(ctx, &editor, false); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
	}

	c.JSON(http.StatusAccepted,
	gin.H{"success": true, "payload": updated, "message": "Job Updated"})
}

func getJobByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	id := c.Param("id") //returns empty string if not there
	if len(id) <= 0 {
		c.JSON(http.StatusBadRequest,
		gin.H{"success": false, "message": "Provide an id"})
		c.Abort(); return
	}

	job, err := jobByID(ctx, id); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
	}

	c.JSON(http.StatusOK,
		gin.H{"success": true, "payload": job})
}

