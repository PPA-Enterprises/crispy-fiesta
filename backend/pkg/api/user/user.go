package user

import (
	"github.com/gin-gonic/gin"
	"PPA"
)


func (self User) Create(c *gin.Context, req PPA.User) {
	/*if err := self.rbac.AccountCreate(c); err != nil {
		return PPA.User{}, err
	}*/
	req.Password = self.securer.Hash(req.Password)
	return self.udb.Create(self.db, req)
}

func (self User) List() {}
func (self User) View() {}
func (self User) Delete() {}

type Update struct {}
func Update() {}
