package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// first there are checking that must be done on incoming request
	err := checkHttpMethod(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	err = checkDefinedPaths(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	proxyReq, err := modifyRequestForProxy(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := sendProxyRequest(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w, err = writeToResponseWriter(resp, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	serverPort, err := getEnv("ServerPort")
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr:    serverPort,
		Handler: http.HandlerFunc(handleRequest),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error %v while serving the requests", err)
	}
}
