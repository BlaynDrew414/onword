package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/programmingbunny/epub-backend/db"
	"github.com/programmingbunny/epub-backend/models"
	"github.com/programmingbunny/epub-backend/responses"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func CreateBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var book models.Book
		defer cancel()

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&book); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newBook := models.Book{
			Title:    book.Title,
			Subtitle: book.Subtitle,
			Author:   book.Author,
		}

		result, err := db.InsertBook(ctx, newBook)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.BookResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}

func GetABook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		bookId := params["bookId"]
		var book models.Book
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(bookId)

		err := db.BookCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&book)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(book)
	}
}

func CreateChapter() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var chapter models.Chapter
		defer cancel()

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&chapter); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&chapter); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newChapter := models.Chapter{
			Title:      chapter.Title,
			ChapterNum: chapter.ChapterNum,
			Text:       chapter.Text,
			BookID:     chapter.BookID,
		}

		result, err := db.InsertChapter(ctx, newChapter)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.BookResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}
