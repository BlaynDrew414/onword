package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chapter struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ImageLocation string             `json:"imageLocation,omitempty" bson:"imageLocation,omitempty"`
	ChapterNum    int                `json:"chapterNum,omitempty" bson:"chapterNum,omitempty"`
	Title         string             `json:"title,omitempty" validate:"required"`
	Text          string             `json:"text,omitempty"`
	BookID        primitive.ObjectID `json:"bookID,omitempty" bson:"bookID,omitempty"`
}

type Chapters struct {
	BookID   primitive.ObjectID `json:"bookID,omitempty" bson:"bookID,omitempty"`
	Chapters []Chapter          `json:"chapters" bson:"chapters"`
}

type ChapterImages struct {
	BookID        primitive.ObjectID `json:"bookID,omitempty" bson:"bookID,omitempty"`
	ChapterNum    int                `json:"chapterNum,omitempty" bson:"chapterNum,omitempty"`
	ImageLocation string             `json:"imageLocation,omitempty" bson:"imageLocation,omitempty"`
	Type          string             `json:"type,omitempty" bson:"type,omitempty"`
}
