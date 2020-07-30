package clients

import (
	"context"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getClientByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	id := c.Param("id") //returns empty string if not there
	if len(id) > 0 {
		c.JSON(http.StatusBadRequest,
		gin.H{"success": false, "message": "Provide an id"})
		c.Abort()
		return
	}

	client, err := clientByID(id, ctx); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"success": true, "payload": client})

}
