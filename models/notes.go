package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notes struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Text      string             `json:"text,omitempty" bson:"text,omitempty"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	BookID    primitive.ObjectID `json:"bookID,omitempty" bson:"bookID,omitempty"`
	VersionID primitive.ObjectID `json:"versionID,omitempty" bson:"versionID,omitempty"`
}
