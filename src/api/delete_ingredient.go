package api

import (
	"bansheechef-server/src/database"
	"net/http"
	"strconv"
)

func DeleteIngredient(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))

	if err != nil {
		_, ok := err.(*strconv.NumError)
		if ok {
			http.Error(response, "expected id", 400)
			return
		}
		
		HandleFatal(&response, err)
		return
	}

	_, err = database.Exec(
		`DELETE FROM ingredients WHERE id = ?;`,
		database.CreateArray(id),
	)

	if err != nil {
		HandleFatal(&response, err)
		return
	}
}
