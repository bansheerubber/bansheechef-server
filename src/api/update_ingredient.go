package api

import (
	"bansheechef-server/src/database"
	"encoding/json"
	"net/http"
	"strconv"
)

//go:generate database-type UpdateIngredientQueryResult
type UpdateIngredientQueryResult struct {
	Name					string	`json:"name"`
	MaxAmount			float64	`json:"maxAmount"`
	ImageSource		string	`json:"image"`
	CurrentAmount	float64	`json:"amount"`
	TypeId				int64		`json:"typeId"`
	Id						int64		`json:"id"`
	Barcode				string	`json:"barcode"`
}

func UpdateIngredient(response http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(1_000)
	if err != nil {
		http.Error(response, "could not parse form", 400)
	}

	ingredientId, err := strconv.ParseInt(request.PostForm.Get("id"), 10, 64)
	if err != nil { // handle id error
		HandleFatal(&response, err)
		return
	}

	amount, err := strconv.ParseFloat(request.PostForm.Get("amount"), 64)
	if err != nil { // handle amount error
		HandleFatal(&response, err)
		return
	}

	result, err := database.Exec(
		`UPDATE ingredients SET current_amount = ? WHERE id = ?;`,
		database.CreateArray(amount, ingredientId),
	)
	if err != nil { // handle update errors
		HandleFatal(&response, err)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 { // if we couldn't update the ingredient, 400 error
		http.Error(response, "no ingredient", 400)
		return
	}

	ingredient, _ := database.QueryOne(
		`SELECT name, max_amount, source, current_amount, i.id, ing.id, barcode
		FROM ingredient_types i
		LEFT JOIN images im ON i.image_id = im.id
		JOIN ingredients ing ON i.id = ing.ingredient_type_id
		WHERE ing.id = ?;`,
		database.CreateArray(ingredientId),
		UpdateIngredientQueryResult_type(),
	)

	binary, err := json.Marshal(ingredient.(*UpdateIngredientQueryResult))
	if err != nil { // handle json error
		HandleFatal(&response, err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(binary)
}
