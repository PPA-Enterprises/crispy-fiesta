package transport

import (
	"PPA"
	"io"
	"pkg/api/job"

	"github.com/gin-gonic/gin"
)

type STREAM struct {
	service job.StreamService
	stream *PPA.StreamEvent
}

func NewStream(service job.StreamService, stream *PPA.StreamEvent, router *gin.RouterGroup, authMw gin.HandlerFunc) {
	httpTransport := STREAM{service:service, stream: stream}
	routes := router.Group("/jobs")
	routes.GET("/", /*authMw,*/ httpTransport.subscribe)
}

func (h STREAM) subscribe(c *gin.Context) {
	//go h.service.Subscribe(c, h.stream)
	go h.service.Subscribe(c, h.stream)
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-h.stream.Message; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
