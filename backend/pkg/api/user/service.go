package user

import(
	"github.com/gin-gonic/gin"
)

type Service interface {
	//Create(gin.Context)
	//List()
	//View()
//	Delete()
	//Update
}

func New() {}
func Init() {}

type User struct {}
type Securer interface {}
type Repository interface {}
type RBAC interface {}
