package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Define the data structure of the API response

// Info is a struct that contains information about mev
type Info struct {
	Name  string   `json:"name"`
	Value string   `json:"value"`
	Alias []string `json:"alias"`
}

// Experience is a struct that contains information about my experience
type Experience struct {
	Name        string   `json:"name"`
	Institution string   `json:"institution"`
	Place       string   `json:"place"`
	Description []string `json:"description"`
	Start       string   `json:"start"`
	End         string   `json:"end"`
}

// Language is a struct that contains information about my languages
type Language struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

// CV is a struct that contains all the information about me
type CV struct {
	Info         []Info       `json:"info"`
	Education    []Experience `json:"education"`
	Experience   []Experience `json:"experience"`
	Technologies []string     `json:"technologies"`
	Tools        []string     `json:"tools"`
	Hobbies      []string     `json:"hobbies"`
	Languages    []Language   `json:"languages"`
	Skills       []string     `json:"skills"`
}

// jsonUnmarshaler is an interface that allows decoding JSON into a struct
type jsonUnmarshaler interface {
	DecodeJSON([]byte) error
}

// DecodeJSON implements the jsonUnmarshaler interface for the CV struct
func (cv *CV) DecodeJSON(data []byte) error {
	return json.Unmarshal(data, cv)
}

// api is a struct that contains the API functions
type api struct{}

// requestAll is a function that requests all the data from the API
func (a api) requestAll() (CV, error) {
	var cv CV
	url := "https://api.lucasvieira.nl/all"

	// Make API request
	if err := a.request(url, &cv); err != nil {
		fmt.Printf("Error getting data from API: %s\n", err)
		return cv, err
	}

	return cv, nil
}

// request is a method that requests data from the API
func (a api) request(url string, j jsonUnmarshaler) error {
	// Create client
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Send request
	res, getErr := client.Do(req)
	if getErr != nil {
		fmt.Println("Error sending request:", getErr)
		return getErr
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Non-ok status code: %d", res.StatusCode)
	}

	// Close response body at end of function
	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}()

	// Read response body
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("Error reading response body:", readErr)
		return readErr
	}

	err = j.DecodeJSON(body)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return err
	}
	return nil
}
