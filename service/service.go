package routes

import (
	"github.com/gorilla/mux"
	books "github.com/programmingbunny/epub-backend/controllers/books"
	chapters "github.com/programmingbunny/epub-backend/controllers/chapters"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/book", books.CreateBook()).Methods("POST")
	router.HandleFunc("/book/{bookId}", books.GetABook()).Methods("GET")

	router.HandleFunc("/chapter", chapters.CreateChapter()).Methods("POST")
	router.HandleFunc("/chapters/{bookId}", chapters.GetAllChapters()).Methods("GET")
}
