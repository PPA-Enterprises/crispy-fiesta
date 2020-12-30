package clients

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func createUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var data createClientCmd
	if c.BindJSON(&data) != nil {
		c.JSON(http.StatusNotAcceptable,
			gin.H{"success": false, "message": "Provide relevant fields"});
		c.Abort(); return
	}
	newClient := fromCreateClientCmd(&data)
	uid, err := newClient.createUniq(ctx); if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()});
		c.Abort(); return
	}
	c.JSON(http.StatusCreated,
		gin.H{"success": true, "payload": uid.String(), "message": "Client Created"})
}

func getClientByPhone(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	phone := c.Param("phone") //returns empty string if not there
	if len(phone) <= 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"success": false, "message": "Provide an id"})
		c.Abort(); return
	}

	client := ClientByPhone(ctx, phone)
	if client == nil {
		c.JSON(http.StatusNotFound,
			gin.H{"success": false, "message": "No client found"})
		c.Abort(); return
	}

	delivarableClient, err := client.Populate(ctx)
	if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
	}

	c.JSON(http.StatusOK,
		gin.H{"success": true, "payload": delivarableClient})

}

func getClientByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	id := c.Param("id") //returns empty string if not there
	if len(id) <= 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"success": false, "message": "Provide an id"})
		c.Abort(); return
	}

	client, err := clientByID(ctx, id)
	if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
	}

	delivarableClient, err := client.Populate(ctx)
	if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
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
		c.Abort(); return
	}

	update, err := tryFromUpdateClientCmd(&data)
	if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
	}

	err = update.Put(ctx, false)
	if err != nil {
		c.JSON(err.Code,
			gin.H{"success": false, "message": err.Error()})
		c.Abort(); return
	}

	c.JSON(http.StatusAccepted,
		gin.H{"success": true, "message": "Job Updated"})

}

func fuzzyClientSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	query := c.Param("query") //returns empty string if not there
	if len(query) <= 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"success": false, "message": "Provide a query"})
		c.Abort(); return
	}

	quantity := c.DefaultQuery("quantity", "100")
	iQuantity, err := strconv.Atoi(quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"success": false, "message": "Invalid Quantity"})
		c.Abort(); return
	}

	results, serachErr := fuzzySearch(ctx, query, iQuantity)
	if serachErr != nil {
		c.JSON(serachErr.Code,
			gin.H{"success": false, "message": serachErr.Error()})
		c.Abort(); return
	}

	c.JSON(http.StatusOK,
		gin.H{"success": true, "payload": populateClients(ctx, results)})
}

// TODO: api/v1/clients?all=bool&sort=bool&source=uint&next=uint
func getClients(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	options := BulkFetchOptions()

	all := c.DefaultQuery("all", "false")
	bAll, err := strconv.ParseBool(all); if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"success": false, "message": "Invalid all Param"})
		c.Abort(); return
	}
	options.All = bAll

	sort := c.DefaultQuery("sort", "false")
	bSort, err := strconv.ParseBool(sort); if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"success": false, "message": "Invalid sort Param"})
		c.Abort(); return
	}
	options.Sort = bSort

	source := c.DefaultQuery("source", "0")
	iSource, err := strconv.ParseUint(source, 10, 64); if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"success": false, "message": "Invalid Source Param"})
		c.Abort(); return
	}
	options.Source = iSource

	next := c.DefaultQuery("next", "10")
	iNext, err := strconv.ParseUint(next, 10, 64); if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"success": false, "message": "Invalid Next Param"})
		c.Abort(); return
	}
	options.Next = iNext

	results, serachErr := fetch(ctx, options)
	if serachErr != nil {
		c.JSON(serachErr.Code,
			gin.H{"success": false, "message": serachErr.Error()})
		c.Abort(); return
	}

	if len(results) > 0 {
		c.JSON(http.StatusOK,
			gin.H{"success": true, "payload": results})
		return
	}

	empty := make([]string, 0)
		c.JSON(http.StatusOK,
			gin.H{"success": true, "payload": empty})
}


