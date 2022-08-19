package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chapter struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	ChapterNum int                `json:"chapterNum,omitempty" bson:"chapterNum,omitempty"`
	Title      string             `json:"title,omitempty" validate:"required"`
	Text       string             `json:"text,omitempty"`
	BookID     primitive.ObjectID `json:"bookID,omitempty" bson:"bookID,omitempty"`
}
