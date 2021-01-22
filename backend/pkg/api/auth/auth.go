package auth
import (
	"net/http"
	"time"
	"context"
	"PPA"
	"github.com/gin-gonic/gin"
)

var (
	InvalidCredentials = PPA.NewAppError(http.StatusUnauthorized, "Incorrect Email or Password")
	NotAuthorized = PPA.NewAppError(http.StatusUnauthorized, "Not Authorized")
)

func (a Auth) Authenticate(c *gin.Context, email, pass string) (*PPA.AuthToken, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	user, err := a.udb.FindByEmail(a.db, ctx, email); if err != nil {
		return nil, err
	}

	if !a.securer.HashMatchesPassword(user.Password, pass) {
		return nil, InvalidCredentials
	}

	/*if !user.Active {
		return nil, NotAuthorized
	}*/

	token, err := a.tokenGen.GenerateToken(*user); if err != nil {
		return nil, NotAuthorized
	}

	user.UpdateLastLogin(a.securer.Token(token))
	if err := a.udb.Update(a.db, ctx, user); err != nil {
		return nil, err
	}
	return &PPA.AuthToken{Token: token, RefreshToken: user.Token}, nil
}

func (a Auth) Refresh(c *gin.Context, refreshToken string) (string, error) {
	duration := time.Now().Add(5*time.Second)
	ctx, cancel := context.WithDeadline(c.Request.Context(), duration)
	defer cancel()

	user, err := a.udb.FindByToken(a.db, ctx, refreshToken); if err != nil {
		return "", nil
	}
	return a.tokenGen.GenerateToken(*user)
}
