package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/programmingbunny/epub-backend/db"
	"github.com/programmingbunny/epub-backend/models"
	"github.com/programmingbunny/epub-backend/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func CreateNotes() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var notes models.Notes
		defer cancel()

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&notes); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&notes); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newNotes := models.Notes{
			Title:     notes.Title,
			Text:      notes.Text,
			Type:      notes.Type,
			BookID:    notes.BookID,
			VersionID: notes.VersionID,
		}

		result, err := db.InsertNotes(ctx, newNotes)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}

func GetNotes() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		noteId := params["noteId"]
		query := r.URL.Query()
		match, err := db.ParseQuery(query)
		if err != nil {
			fmt.Println("err while parsing query in GetNotes Func")
		}
		fmt.Println(match)
		var notes models.Notes
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(noteId)

		err = db.NoteCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&notes)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(notes)
	}
}

func GetAllNotes() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		query := r.URL.Query()
		match, err := db.ParseQuery(query)
		if err != nil {
			fmt.Println("err while parsing query in GetNotes Func")
		}

		var notes models.Notes
		defer cancel()

		err = db.NoteCollection.FindOne(ctx, match).Decode(&notes)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(notes)
	}
}

func UpdateNote() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		noteId := params["noteId"]
		var existingNote models.Notes
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(noteId)

		err := db.NoteCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&existingNote)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Update the existing note
		var updatedNote models.Notes
		err = json.NewDecoder(r.Body).Decode(&updatedNote)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		existingNote.Title = updatedNote.Title
		existingNote.Text = updatedNote.Text

		// Save the updated note to the database
		_, err = db.UpdateNoteByID(ctx, objId, existingNote)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Send a success response back to the client
		rw.WriteHeader(http.StatusOK)
		response := responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": existingNote}}
		json.NewEncoder(rw).Encode(response)
	}
}

// Delete a single note
func DeleteNote() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Create a context with a timeout of 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Get the note ID from the request parameters
		params := mux.Vars(r)
		noteId := params["noteId"]
		log.Printf("Note ID: %s\n", noteId)

		// Convert the note ID to an ObjectID
		objID, err := primitive.ObjectIDFromHex(noteId)
		if err != nil {
			log.Printf("Error converting note ID to ObjectID: %s\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "Invalid note ID", Data: nil}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Delete the note from the database
		result, err := db.DeleteNoteByID(ctx, objID)
		if err != nil {
			log.Printf("Error deleting note: %s\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "Error deleting note", Data: nil}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Check if the note was deleted
		if result.DeletedCount == 0 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.Response{Status: http.StatusNotFound, Message: "Note not found", Data: nil}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Return a success response
		rw.WriteHeader(http.StatusOK)
		response := responses.Response{Status: http.StatusOK, Message: "Note deleted", Data: nil}
		json.NewEncoder(rw).Encode(response)
	}
}
