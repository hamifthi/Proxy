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

func checkHttpMethod(r *http.Request) error {
	if r.Method != http.MethodGet {
		log.Printf("%s method is not supported", r.Method)
		return errors.New(fmt.Sprintf("%s method is not supported", r.Method))
	}
	return nil
}

func checkDefinedPaths(r *http.Request) error {
	rootUrl, err := getEnv("RootUrl")
	if err != nil {
		log.Fatal(err)
	}
	mainPaths := getDefinedPaths(rootUrl)
	targetUrl := r.URL.String()
	parsedUrl, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(parsedUrl.Path, "/")
	if len(parts) < 3 && targetUrl == "/" {
		return nil
	}
	for _, value := range mainPaths {
		if strings.Contains(value, parts[1]) {
			return nil
		}
	}
	log.Printf("Error %s checking destination path", err)
	return errors.New(fmt.Sprintf("Invalid path requested, %s is not defined", parsedUrl.Path))
}
