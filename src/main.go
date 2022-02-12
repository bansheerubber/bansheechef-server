package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"bansheechef-server/src/database"
	"bansheechef-server/src/pages"
)

func main() {
	router := mux.NewRouter()

	// handle index
	router.HandleFunc("/", pages.Index).
		Methods("GET")

	database.Open()
	defer database.Close()

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Serving HTTP on 0.0.0.0:5001")

	log.Fatal(http.ListenAndServe("0.0.0.0:5001", nil))
}
