package auth
import (
	"net/http"
	"PPA"
	"github.com/gin-gonic/gin"
)

func (a Auth) Authenticate(c *gin.Context, user, pass string) (PPA.AuthToken, error) {}
