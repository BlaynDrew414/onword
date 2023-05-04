package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email" validate:"required,email"`
	FirstName string             `bson:"firstName" validate:"required"`
	LastName  string             `bson:"lastName" validate:"required"`
	Password  string             `bson:"password" validate:"required"`
}