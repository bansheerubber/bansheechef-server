package main

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	"bansheechef-server/src/database"
	"bansheechef-server/src/pages"
)

type IngredientType struct {
	Id int
	Name string
	Barcode string
	MaxAmount float32
	IsVolume bool
	UnitCount int
	ImageId int
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", pages.Index).
		Methods("GET")

	database.Open()
	defer database.Close()

	for i := range database.Query("select * from ingredient_types;", nil, reflect.TypeOf((*IngredientType)(nil)).Elem()) {
		log.Println(i.(*IngredientType).Name)
	}

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", nil))
}
