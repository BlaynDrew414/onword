package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Version struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type   string             `json:"type,omitempty" bson:"type"`
	BookID primitive.ObjectID `json:"bookID,omitempty" bson:"bookID,omitempty"`
}
