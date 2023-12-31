package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var rootUrl = "https://swapi.dev/api/"

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// first there are checking that must be done on incoming request
	log.Println("Checking http method")
	err := checkHttpMethod(r)
	if err != nil {
		log.Printf("Error %s checking the http method", err)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	log.Println("Checking url destination")
	err = checkDefinedPaths(r)
	if err != nil {
		log.Printf("Error %s checking destination path", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// creating the new proxy request
	urlParts := strings.Split(r.URL.String(), "/")
	destinationUrl := fmt.Sprintf("%s%s", rootUrl, strings.Join(urlParts[1:3], "/"))
	proxyReq, err := http.NewRequest(http.MethodGet, destinationUrl, r.Body)
	log.Printf("Forwarding request to the %s", destinationUrl)
	if err != nil {
		log.Printf("Error %s creating the proxy request", err)
		http.Error(w, "Error creating the proxy request", http.StatusInternalServerError)
		return
	}

	// adding headers of the original request to the new headers
	log.Println("Adding the request headers to proxy request")
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}
	// set host header
	proxyReq.Host = "Localhost"

	// Create transport and send the request
	customTransport := http.DefaultTransport
	resp, err := customTransport.RoundTrip(proxyReq)
	log.Println("Sending proxy request")
	if err != nil {
		log.Printf("Error %s sending the proxy request", err)
		http.Error(w, "Error sending the proxy request", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// add response headers to the return response to client
	log.Println("Adding the response headers to return response")
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// add http status code to return response
	w.WriteHeader(http.StatusOK)

	// copy the response body to the return response
	log.Println("Writing response body to the return response")
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error %s writing body to the return response", err)
		http.Error(w, "Error writing body to the return response", http.StatusInternalServerError)
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
