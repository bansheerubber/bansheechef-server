package pages

import (
	"net/http"
)

func Index(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "./templates/index.html")
}
