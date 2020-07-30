package clients

import (
	"context"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getClientByPhone(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	phone := c.Param("phone") //returns empty string if not there
	if len(phone) <= 0 {
		c.JSON(http.StatusBadRequest,
		gin.H{"success": false, "message": "Provide an id"})
		c.Abort()
		return
	}

	client := ClientByPhone(ctx, phone); if client == nil {
		c.JSON(http.StatusNotFound,
			gin.H{"success": false, "message": "No client found"})
		c.Abort()
		return
	}

	delivarableClient, err := client.Populate(ctx); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"success": true, "payload": delivarableClient})

}

func update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data updateClientCmd
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable,
		gin.H{"success": false, "message": "Provide relevant fields"})
		c.Abort()
		return
	}

	update, err := tryFromUpdateClientCmd(&data); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	err = update.Put(ctx, false); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusAccepted,
		gin.H{"success": true, "message": "Job Updated"})

}
