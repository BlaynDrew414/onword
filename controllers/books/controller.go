package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("bookPic")
		if err != nil {
			fmt.Println("error while getting the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		tempFile, err := ioutil.TempFile("cover-images", "upload-*.png")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		tempFile.Write(fileBytes)

		book := models.Book{
			Title:     r.FormValue("title"),
			Subtitle:  r.FormValue("subtitle"),
			Author:    r.FormValue("author"),
			BookCover: forwardSlash(tempFile.Name()),
		}

		insertResult, err := db.InsertBook(context.TODO(), book)
		if err != nil {
			fmt.Println(err)
			return
		}
		json.NewEncoder(rw).Encode(insertResult.InsertedID) // return the //mongodb ID of generated document
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

func CreateChapterHeader() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("imageLocation")
		if err != nil {
			fmt.Println("error while getting the File")
			fmt.Println(err)
			return
		}
		defer file.Close()
		tempFile, err := ioutil.TempFile("../chapter-images", "upload-*.png")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		tempFile.Write(fileBytes)

		pass := r.FormValue("bookID")

		newImage := models.ChapterImages{
			BookID:        stringToPrimitive(r.FormValue("bookID")),
			ChapterNum:    stringToInt(r.FormValue("chapterNum")),
			ImageLocation: forwardSlash(tempFile.Name()),
			Type:          r.FormValue("type"),
		}

		insertResult, err := db.InsertImage(context.TODO(), newImage)
		if err != nil {
			log.Fatal(err)
		}

		err = db.UpdateChapterWithHeaderImage(newImage.ImageLocation, pass, newImage.ChapterNum)
		if err != nil {
			fmt.Println(err)
		}
		json.NewEncoder(rw).Encode(insertResult.InsertedID) // return the //mongodb ID of generated document
	}
}

func GetChapterHeader() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		bookId := params["bookId"]
		chNum := params["chNum"]
		var imageLoc models.ChapterImages
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(bookId)
		chapterNum, _ := strconv.Atoi(chNum)

		err := db.ImageCollection.FindOne(ctx, bson.M{"bookID": objId, "chapterNum": chapterNum}).Decode(&imageLoc)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(imageLoc)
	}
}

func forwardSlash(pathName string) string {
	replace := strings.ReplaceAll(pathName, `\`, `/`)
	return replace
}

func stringToInt(input string) int {
	changed, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}
	return changed
}

func stringToPrimitive(input string) primitive.ObjectID {
	objId, _ := primitive.ObjectIDFromHex(input)
	return objId
}
