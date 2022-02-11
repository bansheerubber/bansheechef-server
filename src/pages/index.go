package pages

import (
	"fmt"
	"net/http"
)

func Index(response http.ResponseWriter, request *http.Request) {
	fmt.Println("hey")
}
