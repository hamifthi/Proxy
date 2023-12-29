package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getDefinedPaths(rootUrl string) []string {
	resp, err := http.Get(rootUrl)
	if err != nil {
		log.Fatalf("Error %v while trying to get defined paths", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error %v reading the response of get main paths request", err)
	}
	var mainPathsMap map[string]string
	if err := json.Unmarshal(body, &mainPathsMap); err != nil {
		log.Fatalf("Error %v parsing the response of get main paths request", err)
	}
	var mainPaths []string
	for _, v := range mainPathsMap {
		mainPaths = append(mainPaths, v)
	}
	return mainPaths
}

func checkHttpScheme(r *http.Request) error {
	parsedUrl, err := url.Parse(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	if parsedUrl.Scheme != "http" || parsedUrl.Scheme != "https" {
		return errors.New(fmt.Sprintf("%s scheme is not supported", parsedUrl.Scheme))
	}
	return nil
}

func checkHttpMethod(r *http.Request) error {
	if r.Method != http.MethodGet {
		return errors.New(fmt.Sprintf("%s method is not supported", r.Method))
	}
	return nil
}

func checkDefinedPaths(r *http.Request) error {
	mainPaths := getDefinedPaths(rootUrl)
	targetUrl := r.URL.String()
	parsedUrl, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	for index, value := range mainPaths {
		if !strings.Contains(value, parsedUrl.Path) && index == len(mainPaths)-1 {
			return errors.New(fmt.Sprintf("Invalid path requested, %s is not defined", parsedUrl.Path))
		}
	}
	return nil
}
