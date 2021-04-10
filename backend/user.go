package PPA

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty" m:"Name"`
	Email string `json:"email" bson:"email,omitempty" m:"Email"`
	Password string `json:"-" bson:"password,omitempty"`

	LastLogin time.Time `json:"last_login,omitempty"`
	Token string `json:"-"`
	//Role *Role `json:"role,omitempty" m:"Permission Level"`
	//RoleID AccessRole `json:"-" m:"Access Role"`
	Role string `json:"-" bson:"role,omitempty" m:"Access Level"`
	History []LogEvent `json:"-" bson:"history,omitempty"`
}

type AuthUser struct {
	ID string
	Email string
	Role string
}

func (u *User) UpdateLastLogin(token string) {
	u.Token = token
	u.LastLogin = time.Now()
}

func (u *User) AppendLog(event LogEvent) {
	u.History = append(u.History, event)
}
