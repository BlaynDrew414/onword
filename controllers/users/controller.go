package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/programmingbunny/epub-backend/db"
	"github.com/programmingbunny/epub-backend/models"
	"github.com/programmingbunny/epub-backend/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUsers returns all users in the UserCollection
func GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	cursor, err := db.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a user by their ID from the UserCollection
func GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := db.UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email address from the UserCollection
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := db.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the UserCollection
func CreateUser(ctx context.Context, newUser models.User) error {
    // hash password before storing in the database
    passwordHash, err := db.HashPassword(newUser.Password)
    if err != nil {
        return err
    }
    newUser.Password = passwordHash

    return db.InsertUser(ctx, newUser)
}

// UpdateUser updates an existing user in the UserCollection
func UpdateUser(ctx context.Context, id primitive.ObjectID, updatedUser models.User) error {
	// Update the user document
	update := bson.M{
		"$set": updatedUser,
	}
	_, err := db.UserCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by their ID from the UserCollection
func DeleteUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		id := params["id"]
		objId, _ := primitive.ObjectIDFromHex(id)
		defer cancel()

		// Delete the user from the database
		err := db.DeleteUserByID(ctx, objId)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Send a success response back to the client
		rw.WriteHeader(http.StatusOK)
		response := responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User deleted successfully"}}
		json.NewEncoder(rw).Encode(response)
	}
}