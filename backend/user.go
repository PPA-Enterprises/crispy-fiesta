package crispy_fiesta
import(
"go.mongodb.org/mongo-driver/bson/primitive"
)
type User struct {
	ID primitive.ObjectID `json:"_id,omitempty"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	IsDeleted bool `json:"is_deleted"`
}
