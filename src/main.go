package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"bansheechef-server/src/database"
	"bansheechef-server/src/database/types"
	"bansheechef-server/src/pages"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", pages.Index).
		Methods("GET")

	database.Open()
	defer database.Close()

	for i := range database.Query("select * from ingredient_types;", nil, types.IngredientType_type()) {
		log.Println(i.(*types.IngredientType).Name)
	}

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", nil))
}
