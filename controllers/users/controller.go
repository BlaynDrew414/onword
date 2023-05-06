package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
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

var validate = validator.New()

// GetSingleUser gets a single user by their ID
func GetUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		id := params["id"]
		objId, _ := primitive.ObjectIDFromHex(id)
		defer cancel()

		// Get the user from the database
		user, err := db.GetUserByID(ctx, objId)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Send the user data back to the client
		rw.WriteHeader(http.StatusOK)
		response := responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(rw).Encode(response)
	}
}

// GetUserByEmail retrieves a user by their email address from the UserCollection
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := db.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func CreateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// hash password before storing in the database
		passwordHash, err := db.HashPassword(user.Password)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		user.Password = passwordHash

		err = validate.Struct(user)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		err = db.InsertUser(r.Context(), user)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		json.NewEncoder(rw).Encode(user)
	}
}

func UpdateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		id := params["id"]
		objId, _ := primitive.ObjectIDFromHex(id)
		defer cancel()

		// Parse the updated user from the request body
		var updatedUser models.User
		err := json.NewDecoder(r.Body).Decode(&updatedUser)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Validate the updated user
		validate := validator.New()
		err = validate.Struct(updatedUser)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Update the user in the database
		err = db.UpdateUserById(ctx, objId, updatedUser)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Send a success response back to the client
		rw.WriteHeader(http.StatusOK)
		response := responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User updated successfully"}}
		json.NewEncoder(rw).Encode(response)
	}
}

// DeleteUser deletes a user by their ID from the UserCollection
func DeleteUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := mux.Vars(r)
		id := params["id"]
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Invalid ID"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Delete the user from the database
		err = db.DeleteUserByID(ctx, objId)
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
