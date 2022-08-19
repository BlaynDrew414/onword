package main

import (
	"log"
	"net/http"

	"github.com/programmingbunny/epub-backend/configs"
	routes "github.com/programmingbunny/epub-backend/service"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(router) //add this

	log.Fatal(http.ListenAndServe(":3000", router))
}
