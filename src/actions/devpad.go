package actions

import (
	"blog-server/database"
	"blog-server/types"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func FetchProjects(userID int, force bool) (*types.ProjectCache, error) {
	// if we have a cached version of the projects (in recent 24 hours) with a status of 'fetched', then we can return that row. otherwise, we must fetch via the api & cache this result, which requires us to INSERT a 'pending' row, make the fetch, then either update the status to 'fetched' or 'failed' depending on the response, and then we return the updated cached row
	var cache types.ProjectCache
	if !force {
		cache, err := database.GetLatestProjectCache(userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
			} else {
				return nil, errors.Join(fmt.Errorf("Error getting latest project cache"), err)
			}
		}
		if cache.Status == "fetched" {
			return &cache, nil
		}
	}

	devpad_url := os.Getenv("DEVPAD_API")

	// check to see that the user has a token available
	key, err := database.GetProjectKey(userID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Error getting project key"), err)
	}
	log.Debug("Fetching projects", "url", devpad_url+"/projects", "key", key)
	// insert 'pending' row
	cache = types.ProjectCache{
		UserID:    userID,
		Status:    "pending",
		URL:       devpad_url + "/projects",
		Data:      "",
		FetchedAt: time.Now(),
	}
	// we will need the ID of the row that we're about to insert
	id, err := database.InsertProjectCache(cache)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Error inserting project cache"), err)
	}
	log.Debug("Inserted project cache", "id", id)

	// make the request to the devpad api
	// put the key in the header "Authorization: Bearer <key>"
	client := &http.Client{}
	req, err := http.NewRequest("GET", devpad_url+"/projects", nil)
	if err != nil {
		database.FailProjectCache(id)
		return nil, errors.Join(fmt.Errorf("Error creating request"), err)
	}
	req.Header = http.Header{
		"Authorization": {"Bearer " + key},
	}
	resp, err := client.Do(req)
	if err != nil {
		database.FailProjectCache(id)
		return nil, errors.Join(fmt.Errorf("Error fetching projects api"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		database.FailProjectCache(id)
		log.Error("Error fetching projects", "status", resp.StatusCode, "url", devpad_url+"/projects")
		return nil, errors.New("Bad status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		database.FailProjectCache(id)
		return nil, errors.Join(fmt.Errorf("Error reading body"), err)
	}

	if db_err := database.UpdateProjectCacheData(id, string(body)); db_err != nil {
		return nil, errors.Join(fmt.Errorf("Error updating project cache data"), db_err)
	}

	cache, err = database.GetProjectCache(id)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Error getting project cache"), err)
	}

	return &cache, nil

}
