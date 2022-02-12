package api

import (
	"bansheechef-server/src/database"
	"encoding/json"
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

	result, err := database.QueryOne(
		`SELECT name, max_amount, source, i.id
		FROM ingredient_types i
		LEFT JOIN images im ON i.image_id = im.id
		WHERE barcode = ?;`,
		database.CreateArray(barcode),
		BarcodeQueryResult_type(),
	)

	if err != nil {
		HandleFatal(&response, err)
		return
	}

	if result == nil { // send empty json if we couldn't get a result
		response.Header().Set("Content-Type", "application/json")
		response.Write([]byte("{}"))
		return
	}

	binary, err := json.Marshal(result.(*BarcodeQueryResult))
	if err != nil {
		HandleFatal(&response, err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(binary)
}
