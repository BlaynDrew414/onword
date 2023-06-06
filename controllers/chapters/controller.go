package chapters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/programmingbunny/epub-backend/db"
	"github.com/programmingbunny/epub-backend/models"
	"github.com/programmingbunny/epub-backend/responses"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func CreateChapter() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var chapter models.Chapter
		defer cancel()

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&chapter); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&chapter); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newChapter := models.Chapter{
			Title:      chapter.Title,
			ChapterNum: chapter.ChapterNum,
			Text:       chapter.Text,
			BookID:     chapter.BookID,
			VersionID:  chapter.VersionID,
		}

		nums, err := getBookNumbers(newChapter.BookID)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := range nums {
			if newChapter.ChapterNum == nums[i] {
				newChapter.ChapterNum = newChapter.ChapterNum + 1
			}
		}

		result, err := db.InsertChapter(ctx, newChapter)
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

func GetAllChapters() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		bookId := params["bookId"]
		var chapters []models.Chapter
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(bookId)

		results, err := db.ChapterCollection.Find(ctx, bson.M{"bookID": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleChapter models.Chapter
			if err = results.Decode(&singleChapter); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
			}

			chapters = append(chapters, singleChapter)
		}

		var everyChapter models.Chapters
		everyChapter.Chapters = append(everyChapter.Chapters, chapters...)
		everyChapter.BookID = objId

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(everyChapter)
	}
}

func GetSingleChapter() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		chapterId := params["chapterId"]
		var chapter models.Chapter
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(chapterId)

		err := db.ChapterCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&chapter)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(chapter)
	}
}

func UpdateChapter() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		chapterId := params["chapterId"]
		var existingChapter models.Chapter
		defer cancel()

		objId, err := primitive.ObjectIDFromHex(chapterId)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		err = db.ChapterCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&existingChapter)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Update the existing chapter
		var updatedChapter models.Chapter
		err = json.NewDecoder(r.Body).Decode(&updatedChapter)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		existingChapter.Title = updatedChapter.Title
		existingChapter.Text = updatedChapter.Text
		existingChapter.ChapterNum = updatedChapter.ChapterNum

		// Save the updated chapter to the database
		_, err = db.UpdateChapterByID(ctx, objId, existingChapter)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Send a success response back to the client
		rw.WriteHeader(http.StatusOK)
		response := responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": existingChapter}}
		json.NewEncoder(rw).Encode(response)
	}
}

func getBookNumbers(bookId primitive.ObjectID) ([]int, error) {
	var getNumbers []int

	results, err := db.ChapterCollection.Find(context.TODO(), bson.M{"bookID": bookId})
	if err != nil {
		return nil, err
	}

	for results.Next(context.TODO()) {
		var singleChapter models.Chapter
		err = results.Decode(&singleChapter)
		if err != nil {
			return nil, err
		}
		getNumbers = append(getNumbers, singleChapter.ChapterNum)
	}

	results.Close(context.TODO())

	return getNumbers, nil
}

func DeleteChapter() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := mux.Vars(r)
		chapterId := params["chapterId"]

		objId, err := primitive.ObjectIDFromHex(chapterId)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		result, err := db.DeleteChapterByID(ctx, objId)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount == 0 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "chapter not found"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "chapter deleted successfully"}}
		json.NewEncoder(rw).Encode(response)
	}
}
