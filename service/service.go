package routes

import (
	"github.com/programmingbunny/epub-backend/controllers"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/book", controllers.CreateBook()).Methods("POST")
	router.HandleFunc("/book/{bookId}", controllers.GetABook()).Methods("GET")

	router.HandleFunc("/chapter", controllers.CreateChapter()).Methods("POST")
}
