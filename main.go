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

	client := configs.DB

	//routes
	routes.Routes(router, client) 

	log.Println("Hello, This is OnWord!")
	log.Fatal(http.ListenAndServe(":3000", router))
}
