package api

import (
	"bansheechef-server/src/database"
	"database/sql"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

//go:generate database-type AddIngredientQueryResult
type AddIngredientQueryResult struct {
	TypeId			int64
	ImageId			sql.NullInt64
	ImageSource sql.NullString
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")
func getRandomImageName() string {
	randomName := make([]rune, 8)
	for i := range randomName {
		randomName[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomName)
}

func AddIngredient(response http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(10_000_000)
	if err != nil {
		http.Error(response, "could not parse form", 400)
	}
	
	name := request.PostForm.Get("name")
	maxAmount, maxAmountErr := strconv.ParseFloat(request.PostForm.Get("maxAmount"), 64)
	barcode := request.PostForm.Get("barcode")
	units := request.PostForm.Get("units")

	// check if things are null or not
	if name == "" || units == "" {
		http.Error(response, "expected name, max amount, units", 400)
		return
	}
	
	if maxAmountErr != nil {
		HandleFatal(&response, maxAmountErr)
		return
	}

	currentAmount := maxAmount
	if request.PostForm.Get("currentAmount") != "" { // default to using maxAmount if we could not convert the current amount
		currentAmount2, currentAmountErr := strconv.ParseFloat(request.PostForm.Get("currentAmount"), 64)
		currentAmount = currentAmount2
		if currentAmountErr != nil {
			HandleFatal(&response, currentAmountErr)
			return
		}
	}

	// check if we already have an ingredient type
	result, err := database.QueryOne(
		`SELECT i.id, im.id, im.source
		FROM ingredient_types i
		LEFT JOIN images im on im.id = i.image_id
		WHERE name = ? AND max_amount = ?;`,
		database.CreateArray(name, maxAmount),
		AddIngredientQueryResult_type(),
	)

	if err != nil {
		HandleFatal(&response, err)
		return
	}

	// handle image stuff
	var ingredientTypeId int64 = 0
	var imageId int64 = 0
	image := ""
	if result != nil { // if we already have a type in the database, just use that
		ingredientType := *result.(*AddIngredientQueryResult)
		ingredientTypeId = ingredientType.TypeId
		imageId = ingredientType.ImageId.Int64
		image = ingredientType.ImageSource.String
	} else {
		imageName := getRandomImageName()
		image = "local:" + imageName
		imageFileName := filepath.Join(database.LOCAL_IMAGES, imageName)

		file, _, err := request.FormFile("pictureBlob")
		if file != nil { // if the client gave us file data
			imageFile, _ := os.Create(imageFileName)
			defer imageFile.Close()
			_, err = io.Copy(imageFile, file)

			if err != nil {
				HandleFatal(&response, err)
				return
			}
			
			result, _ := database.Exec(
				`INSERT INTO images (source) VALUES(?);`,
				database.CreateArray(image),
			)
			imageId, _ = result.LastInsertId()
		} else if request.PostForm.Get("picture") != "" { // if the client gave us a picture url
			picture := request.PostForm.Get("picture")
			
			// download the picture
			pictureResponse, err := http.Get(picture)
			if err != nil {
				HandleFatal(&response, err)
				return
			}

			defer pictureResponse.Body.Close()

			if pictureResponse.StatusCode != 200 {
				http.Error(response, "invalid picture URL reqeuest response", 400)
				return
			}

			// write the image we downloaded to file
			imageFile, _ := os.Create(imageFileName)
			defer imageFile.Close()
			_, err = io.Copy(imageFile, pictureResponse.Body)
			if err != nil {
				HandleFatal(&response, err)
				return
			}

			result, _ := database.Exec(
				`INSERT INTO images (source) VALUES(?);`,
				database.CreateArray(image),
			)
			imageId, _ = result.LastInsertId()
		}

		// insert ingredient type into database
		result, err := database.Exec(
			`INSERT INTO ingredient_types (name, max_amount, image_id, unit_count, is_volume, barcode)
			VALUES(?, ?, ?, 0, TRUE, ?);`,
			database.CreateArray(name, maxAmount, imageId, barcode),
		)
		if err != nil {
			HandleFatal(&response, err)
			return
		}

		// get the ingredient type id
		ingredientTypeId, _ = result.LastInsertId()
	}

	// add the ingredient
	ingredientResult, _ := database.Exec(
		`INSERT INTO ingredients (ingredient_type_id, current_amount) VALUES(?, ?);`,
		database.CreateArray(ingredientTypeId, currentAmount),
	)
	ingredientId, _ := ingredientResult.LastInsertId()

	// convert to JSON
	jsonResult := struct {
		Amount		float64 `json:"amount"`
		Barcode		string	`json:"barcode"`
		Id				int64		`json:"id"`
		Image			string	`json:"image"`
		MaxAmount	float64	`json:"maxAmount"`
		Name			string	`json:"name"`
		TypeId		int64		`json:"typeId"`
	} {
		currentAmount,
		barcode,
		ingredientId,
		image,
		maxAmount,
		name,
		ingredientTypeId,
	};
	binary, err := json.Marshal(jsonResult)
	if err != nil {
		HandleFatal(&response, err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(binary)
}
