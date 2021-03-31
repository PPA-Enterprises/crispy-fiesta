package PPA
import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

type AppError struct {
	Status int `json:"-"`
	Success bool `json:"success"`
	Message string `json:"message,omitempty"`
}

var (
	InternalError = ErrWithStatus(http.StatusInternalServerError)
	DbError = ErrWithStatus(http.StatusInternalServerError)
	AlreadyExists = ErrWithStatus(http.StatusConflict)
	Forbidden = ErrWithStatus(http.StatusForbidden)
	BadRequest = ErrWithStatus(http.StatusBadRequest)
	NotFound = ErrWithStatus(http.StatusNotFound)
	Unauthorized = ErrWithStatus(http.StatusUnauthorized)
)

var validationErrors = map[string]string {
	"required": " is required, but not received",
	"email": " is not a valid email",
}

func ErrWithStatus(status int) *AppError {
	return &AppError{Status: status, Success: false}
}

func NewAppError(status int, msg string) *AppError {
	return &AppError{Status: status, Success: false, Message: msg}
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
	if err != nil {
		fmt.Println(err)
	}
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
			msg = append(msg, fmt.Sprintf("%s%s", v.StructField(), validationErrorMessage(v.ActualTag())))
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": msg})

	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H {
			"success": false,
			"message": err.Error(),
		})
	}
}
