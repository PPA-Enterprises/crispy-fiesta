package app_error
import (
	"http"
	"fmt"
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v10"
)

type AppError struct {
	Status int `json:"-"`
	Message string `json:"message,omitempty"`
}

var (
	Internal = NewStatus(http.StatusInternalServerError)
	DB = NewStatus(http.StatusInternalServerError)
	Forbidden = NewStatus(http.StatusForbidden)
	BadRequest = NewStatus(http.StatusBadRequest)
	NotFound = NewStatus(http.StatusNotFound)
	Unauthorized = NewStatus(http.StatusUnauthorized)
)

var validationErrors = map[string]string {
	"required": " is required, but not received",
	"email": " is not a valid email",
}

func NewStatus(status int) *AppError {
	return &AppError{Status: status}
}

func New(status int, msg string) *AppError {
	return &AppError{Status: status, Message: msg}
}

// Needed to implement Go error interface
func (e AppError) Error() string {
	return e.Message
}

func validationErrorMessage(s string) string {
	if v, ok := validationErrors[s]; ok {
		return v
	}
	return " failed on " + s + " validation"
}

func Response(c *gin.Context, err error) {
	switch err.(type) {
	case *AppError:
		e := err.(*AppError)
		if e.Message == "" {
			c.AbortWithStatusJSON(e.Status, gin.H{"success": false})
		} else {
			c.AbortWithStatusJSON(e.Status, gin.H{"success": false, "message": e.Message})
		}
		return

	case validator.ValidationErrors:
		var msg []string
		e := err.(validator.ValidationErrors)

		for _, v := range e {
			msg = append(msg, fmt.Sprintf("%s%s", v.Name, validationErrorMessage(v.ActualTag)))
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": msg})

	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H {
			"success": false,
			"message": err.Error(),
		})
	}
}
