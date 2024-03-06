package actions

import (
	"blog-server/database"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

// this file is responsible for taking a queue_link (integration) and fetching data from API and syncing posts to database

func SyncUserDevTo(userID int) error {
	// get all integrations for the user
	link, err := database.GetIntegration(userID, "devto")
	if err != nil {
		return err
	}
	if link == nil {
		return errors.New("No integration found")
	}

	// fetch data from the location
	url := link.Location
	// we need to get the token out of the 'data' json field
	var data map[string]interface{}
	err = json.Unmarshal([]byte(link.Data), &data)

	if err != nil {
		return err
	}
	token, ok := data["token"].(string)
	if !ok {
		return errors.New("No token found in data")
	}

    log.Info("Fetching DevTo API", "url", url, "token", token)

	result, err := fetchDevToAPI(url, token)
	if err != nil {
		return err
	}

    // do something with result, for now just print out to console
    log.Info("DevTo API resul", "result", result)

	return nil

}

func fetchDevToAPI(url string, token string) (interface{}, error) {
	// make http request to url with token in the headers
	// location should be https://dev.to/api/articles/me
	// we want to make a GET request with 'api-key' set to token and 'accept' as 'application/vnd.forem.api-v1+json'
    
    // Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set request headers
	req.Header.Set("api-key", token)
	req.Header.Set("accept", "application/vnd.forem.api-v1+json")

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Decode the response body
	var responseData interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
