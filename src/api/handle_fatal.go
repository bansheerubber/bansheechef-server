package api

import (
	"log"
	"net/http"
)

func HandleFatal(response *http.ResponseWriter, err error) {
	http.Error(*response, "fatal error", 502)
	log.Print(err)
}
