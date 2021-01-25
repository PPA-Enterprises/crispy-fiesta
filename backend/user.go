package PPA

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`
	Password string `json:"-" bson:"password,omitempty"`

	LastLogin time.Time `json:"last_login,omitempty"`
	Token string `json:"-"`
	Role *Role `json:"role,omitempty"`
	RoleID AccessRole `json:"-"`
}

type AuthUser struct {
	ID string
	Email string
	Role AccessRole
}

func (u *User) UpdateLastLogin(token string) {
	u.Token = token
	u.LastLogin = time.Now()
}
