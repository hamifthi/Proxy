package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func modifyRequestForProxy(r *http.Request) (*http.Request, error) {
	urlParts := strings.Split(r.URL.String(), "/")
	destinationUrl := fmt.Sprintf("%s%s", rootUrl, strings.Join(urlParts[1:], "/"))
	proxyReq, err := http.NewRequest(http.MethodGet, destinationUrl, r.Body)

	if err != nil {
		log.Printf("Error %s creating the proxy request", err)
		return nil, fmt.Errorf("error %s creating the proxy request", err)
	}

	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	proxyReq.Host = "Localhost"

	return proxyReq, nil
}

func sendProxyRequest(proxyReq *http.Request) (*http.Response, error) {
	customTransport := http.DefaultTransport
	resp, err := customTransport.RoundTrip(proxyReq)

	if err != nil {
		log.Printf("Error %s sending the proxy request", err)
		return nil, fmt.Errorf("error %s sending the proxy request", err)
	}

	return resp, nil
}

func writeToResponseWriter(resp *http.Response, w http.ResponseWriter) (http.ResponseWriter, error) {
	// add response headers to the return response to client
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// add http status code to return response
	w.WriteHeader(resp.StatusCode)

	// copy the response body to the return response
	log.Println("Writing response body to the return response")
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error %s writing body to the return response", err)
		return nil, fmt.Errorf("error %s writing proxy response body to main response", err)
	}
	defer resp.Body.Close()

	return w, nil
}
