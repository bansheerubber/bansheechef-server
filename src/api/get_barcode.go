package api

import (
	"bansheechef-server/src/database"
	"encoding/json"
	"log"
	"net/http"
)

//go:generate database-type BarcodeQueryResult
type BarcodeQueryResult struct {
	Name 				string	`json:"name"`
	MaxAmount 	float32	`json:"maxAmount"`
	ImageSource string	`json:"image"`
	TypeId			int			`json:"typeId"`
}

func GetBarcode(response http.ResponseWriter, request *http.Request) {
	barcode := request.URL.Query().Get("barcode")

	if barcode == "" {
		http.Error(response, "expected barcode", 400)
		return
	}

	result := database.QueryOne(
		`SELECT name, max_amount, source, i.id
		FROM ingredient_types i
		LEFT JOIN images im ON i.image_id = im.id
		WHERE barcode = ?;`,
		database.CreateArray(barcode),
		BarcodeQueryResult_type(),
	).(*BarcodeQueryResult)

	if result == nil { // send empty json if we couldn't get a result
		response.Header().Set("Content-Type", "application/json")
		response.Write([]byte("{}"))
		return
	}

	binary, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(binary)
}
