package secure
import (
	"github.com/gin-gonic/gin"
	"PPA/middleware"
)

func Add(r *gin.Engine, h ...gin.HandlerFunc) {
	for _, v := range h {
		r.Use(v)
	}
}

func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-DNS-Prefetch-Control", "off")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Strict-Transport-Security", "max-age=5184000; includeSubDomains")
		c.Header("X-Download-Options", "noopen")
		c.Header("X-XSS-Protection", "1; mode=block")
	}
}

