package PPA
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Label struct {
	ID primitive.ObjectID `json:"_id",omitempty" bson:"_id,omitempty"`
	LabelName `json:"label_name" bson:"label_name,omitempty"`
	Clients []primitive.ObjectID `json:"clients" bson:"clients,omitempty"`
	History []LogEvent `"json:"history" bson:"history,omitempty"`
}

func( l *Label) AppendLog(event LogEvent) {
	l.History = append(l.History, event)
}
