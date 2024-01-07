package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestInvalidHttpMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleRequest))
	defer ts.Close()

	// Make a POST request to the test server
	res, err := http.Post(ts.URL, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code returned
	if status := res.StatusCode; status != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestUndefinedPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleRequest))
	defer ts.Close()

	// Make a GET request to an undefined path
	res, err := http.Get(fmt.Sprintf("%s/undefined_path", ts.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code returned
	if status := res.StatusCode; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestDefaultPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleRequest))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check if status code is 200 OK
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var definedPathsMap map[string]string
	if err := json.Unmarshal(body, &definedPathsMap); err != nil {
		t.Fatal(err)
	}
	expectedDefinedPathsMap := map[string]string{
		"people":    "https://swapi.dev/api/people/",
		"planets":   "https://swapi.dev/api/planets/",
		"films":     "https://swapi.dev/api/films/",
		"species":   "https://swapi.dev/api/species/",
		"vehicles":  "https://swapi.dev/api/vehicles/",
		"starships": "https://swapi.dev/api/starships/",
	}
	if !reflect.DeepEqual(expectedDefinedPathsMap, definedPathsMap) {
		t.Errorf("Expected %s, got %s", expectedDefinedPathsMap, definedPathsMap)
	}
}

func TestDefinedPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleRequest))
	defer ts.Close()

	// Make a GET request to the 'people' path
	res, err := http.Get(fmt.Sprintf("%s/people/", ts.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code is 200
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Read the body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// We don't know the exact response, but we can check if it includes some expected strings
	// This is just an example. Please replace it with your own logic.
	if string(body) == "" {
		t.Errorf("response does not contain anything got %s", string(body))
	}
}

func TestDefinedPathWithParameters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleRequest))
	defer ts.Close()

	// Make a GET request to the 'people' path
	res, err := http.Get(fmt.Sprintf("%s/people/1", ts.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code is 200
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Read the body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// We don't know the exact response, but we can check if it includes some expected strings
	// This is just an example. Please replace it with your own logic.
	if string(body) == "" {
		t.Errorf("response does not contain anything got %s", string(body))
	}
}
