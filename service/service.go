package routes

import (
	"github.com/gorilla/mux"
	books "github.com/programmingbunny/epub-backend/controllers/books"
	chapters "github.com/programmingbunny/epub-backend/controllers/chapters"
	"github.com/programmingbunny/epub-backend/controllers/notes"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(router *mux.Router, client *mongo.Client) {
	router.HandleFunc("/createBook", books.CreateBook()).Methods("POST")
	router.HandleFunc("/book/{bookId}", books.GetABook()).Methods("GET")
	router.HandleFunc("/deleteBook/{bookId}", books.DeleteBook()).Methods("Delete")

	router.HandleFunc("/chapter", chapters.CreateChapter()).Methods("POST")
	router.HandleFunc("/chapters/{bookId}", chapters.GetAllChapters()).Methods("GET")
	router.HandleFunc("/chapter/{chapterId}", chapters.GetSingleChapter()).Methods("GET")
	//router.HandleFunc("/updateChapterTitle/{bookId}/{chNum}", chapters.UpdateChapterTitle()).Methods("PUT")

	router.HandleFunc("/getChapterImage/{bookId}/{chNum}", books.GetChapterHeader()).Methods("GET")
	router.HandleFunc("/createChapterImage", books.CreateChapterHeader()).Methods("POST")

	router.HandleFunc("/getNotes", notes.GetAllNotes()).Methods("GET")
	router.HandleFunc("/getNotes/{noteId}", notes.GetNotes()).Methods("GET")
	router.HandleFunc("/createNotes", notes.CreateNotes()).Methods("POST")
	router.HandleFunc("/updateNotes/{noteId}", notes.UpdateNote()).Methods("PUT")
	router.HandleFunc("/deleteNotes/{noteId}", notes.DeleteNote()).Methods("DELETE")
	
}
