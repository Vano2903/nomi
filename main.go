package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	//r must handle those function and allowd only get method
	r.HandleFunc("/names", NamesHandler).Methods("GET")
	r.HandleFunc("/name", NameHandler).Methods("GET")

	//this root use /{toSearch}, mux allow us to declare url query even like this
	//to find a name now we need to search /exist/<name> and will handle it normally
	r.HandleFunc("/exist/{toSearch}", ExistHandler).Methods("GET")

	//read port from enviroment, if not found will assing 8080 by default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
