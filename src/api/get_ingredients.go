package api

import (
	"bansheechef-server/src/database"
	"encoding/json"
	"log"
	"net/http"
)

//go:generate database-type GetIngredientsResult
type GetIngredientsResult struct {
	Name 							string	`json:"name"`
	MaxAmount 				float32	`json:"maxAmount"`
	ImageSource 			string	`json:"image"`
	CurrentAmount 		float32	`json:"amount"`
	IngredientTypeId 	int			`json:"typeId"`
	IngredientId 			int			`json:"id"`
	Barcode 					string	`json:"barcode"`
}

func GetIngredients(response http.ResponseWriter, request *http.Request) {
	var results []GetIngredientsResult
	query := database.Query(
		`SELECT name, max_amount, source, current_amount, i.id, ing.id, barcode
		FROM ingredient_types i
		LEFT JOIN images im ON i.image_id = im.id
		JOIN ingredients ing ON i.id = ing.ingredient_type_id;`,
		nil,
		GetIngredientsResult_type(),
	)

	for i := range query {
		results = append(results, *query[i].(*GetIngredientsResult))
	}

	binary, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(binary)
}
