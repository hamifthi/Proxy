package main

import (
	"log"
	"net/http"
)

var rootUrl = "https://swapi.dev/api/"

func handleRequest(w http.ResponseWriter, r *http.Request) {
	err := checkHttpScheme(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUpgradeRequired)
		return
	}
	err = checkHttpMethod(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	err = checkDefinedPaths(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handleRequest),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error %v while serving the requests", err)
	}
}
