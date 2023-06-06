package routes

import (
	"github.com/gorilla/mux"
	books "github.com/programmingbunny/epub-backend/controllers/books"
	chapters "github.com/programmingbunny/epub-backend/controllers/chapters"
	"github.com/programmingbunny/epub-backend/controllers/notes"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/programmingbunny/epub-backend/controllers/users"
	
)

func Routes(router *mux.Router, client *mongo.Client) {
	router.HandleFunc("/createUser", users.CreateUser()).Methods("POST")
	router.HandleFunc("/getUser/{userId}", users.GetUser()).Methods("GET")
	router.HandleFunc("/deleteUser/{userId}", users.DeleteUser()).Methods("DELETE")
	router.HandleFunc("/updateUser/{userId}", users.UpdateUser()).Methods("PUT")


	router.HandleFunc("/createBook", books.CreateBook()).Methods("POST")
	router.HandleFunc("/book/{bookId}", books.GetABook()).Methods("GET")
	router.HandleFunc("/deleteBook/{bookId}", books.DeleteBook()).Methods("Delete")

	router.HandleFunc("/createChapter", chapters.CreateChapter()).Methods("POST")
	router.HandleFunc("/getChapters/{bookId}", chapters.GetAllChapters()).Methods("GET")
	router.HandleFunc("/getChapter/{chapterId}", chapters.GetSingleChapter()).Methods("GET")
	router.HandleFunc("/updateChapter/{bookId}/{chapterId}", chapters.UpdateChapter()).Methods("PUT")
	router.HandleFunc("/deleteChapter/{chapterId}", chapters.DeleteChapter()).Methods("DELETE")

	router.HandleFunc("/getChapterImage/{bookId}/{chapterId}", books.GetChapterHeader()).Methods("GET")
	router.HandleFunc("/createChapterImage", books.CreateChapterHeader()).Methods("POST")

	router.HandleFunc("/getNotes", notes.GetAllNotes()).Methods("GET")
	router.HandleFunc("/getNotes/{noteId}", notes.GetNotes()).Methods("GET")
	router.HandleFunc("/createNotes", notes.CreateNotes()).Methods("POST")
	router.HandleFunc("/updateNotes/{noteId}", notes.UpdateNote()).Methods("PUT")
	router.HandleFunc("/deleteNotes/{noteId}", notes.DeleteNote()).Methods("DELETE")
	
}
