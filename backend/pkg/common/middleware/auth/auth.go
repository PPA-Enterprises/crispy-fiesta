package auth

import(
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}

func Middleware(tokenParser TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := tokenParser.ParseToken(c.Request.Header.Get("Authorization"))
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"success": false, "Message": "Not Authorized"}); return
		}

		claims := token.Claims.(jwt.MapClaims)

		id := string(claims["_id"].(string))
		name := string(claims["name"].(string))

		c.Set("_id", id)
		c.Set("name", name)

		c.Next()
	}
}
