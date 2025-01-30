package actions

import (
	"blog-server/database"
	"blog-server/types"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

func FetchProjects(userID int) (*types.ProjectCache, error) {
	// if we have a cached version of the projects (in recent 24 hours) with a status of 'fetched', then we can return that row. otherwise, we must fetch via the api & cache this result, which requires us to INSERT a 'pending' row, make the fetch, then either update the status to 'fetched' or 'failed' depending on the response, and then we return the updated cached row
	cache, err := database.GetLatestProjectCache(userID)
	if err != nil {
		return nil, err
	}
	if cache.Status == "fetched" {
		return &cache, nil
	}

	devpad_url := os.Getenv("DEVPAD_API_URL")

	// check to see that the user has a token available
	key, err := database.GetProjectKey(userID)
	if err != nil {
		return nil, err
	}
	// insert 'pending' row
	cache = types.ProjectCache{
		UserID:    userID,
		Status:    "pending",
		URL:       "https://devpad.dev/api/projects",
		Data:      sql.NullString{},
		FetchedAt: time.Now(),
	}
	// we will need the ID of the row that we're about to insert
	id, err := database.InsertProjectCache(cache)
	if err != nil {
		return nil, err
	}

	// make the request to the devpad api
	// put the key in the header "Authorization: Bearer <key>"
	client := &http.Client{}
	req, err := http.NewRequest("GET", devpad_url+"/projects", nil)
	if err != nil {
		database.FailProjectCache(id)
		return nil, err
	}
	req.Header = http.Header{
		"Authorization": {"Bearer " + key},
	}
	resp, err := client.Do(req)
	if err != nil {
		database.FailProjectCache(id)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		database.FailProjectCache(id)
		return nil, errors.New("Error fetching projects")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		database.FailProjectCache(id)
		return nil, err
	}

	if db_err := database.UpdateProjectCacheData(id, body); db_err != nil {
		return nil, db_err
	}

	cache, err = database.GetProjectCache(id)
	if err != nil {
		return nil, err
	}

	return &cache, nil

}
