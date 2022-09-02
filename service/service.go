package routes

import (
	"github.com/gorilla/mux"
	books "github.com/programmingbunny/epub-backend/controllers/books"
	chapters "github.com/programmingbunny/epub-backend/controllers/chapters"
)

func Routes(router *mux.Router) {
	router.HandleFunc("/createBook", books.CreateBook()).Methods("POST")
	router.HandleFunc("/book/{bookId}", books.GetABook()).Methods("GET")

	router.HandleFunc("/chapter", chapters.CreateChapter()).Methods("POST")
	router.HandleFunc("/chapters/{bookId}", chapters.GetAllChapters()).Methods("GET")
	router.HandleFunc("/chapter/{chapterId}", chapters.GetSingleChapter()).Methods("GET")

	router.HandleFunc("/getChapterImage/{bookId}/{chNum}", books.GetChapterHeader()).Methods("GET")
	router.HandleFunc("/createChapterImage", books.CreateChapterHeader()).Methods("POST")
}
