package chapters

import (
	"context"
	"encoding/json"
	"fmt"
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
			response := responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.BookResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
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
			response := responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleChapter models.Chapter
			if err = results.Decode(&singleChapter); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
			}

			chapters = append(chapters, singleChapter)
		}

		var everyChapter models.Chapters
		for i := range chapters {
			everyChapter.Chapters = append(everyChapter.Chapters, chapters[i])
		}
		everyChapter.BookID = objId

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(everyChapter)
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
