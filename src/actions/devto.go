package actions

import (
	"blog-server/database"
	"blog-server/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

// this file is responsible for taking a queue_link (integration) and fetching data from API and syncing posts to database

func SyncUserDevTo(userID int) error {
    // first fetch user account
    user, err := database.GetUserByID(userID)
    if err != nil {
        return err
    }
	// get all integrations for the user
	link, err := database.GetIntegrationBySource(userID, "devto")
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

    /**
     * with this data we then want to cross check against the existing posts
     * so for each post we are going to check to see if there is already a hard link in fetch_links
     * fetch_links has column 'identifier' which we will take from the devto api ('slug') and 'fetch_source' which is the id of the integration (link.ID)
     * if there is no match then we want to check against existing posts with the same name
     * if we find a post with the same name then we just create a fetch_link instance linking that post to this integration
     * otherwise we create a new post based on the data from the devto api and then create a fetch_link instance linking that post to this integration
     */

    // iterate over devto articles
    articles, ok := result.([]interface{})
    if !ok {
        return errors.New("Invalid response from DevTo API")
    }

    for _, article := range articles {
        article_data, ok := article.(map[string]interface{})
        if !ok {
            return errors.New("Invalid article data from DevTo API")
        }
        err := handleDevToArticle(article_data, user, link)
        if err != nil {
            log.Error("Error handling DevTo article", "error", err)
            continue
        }
    }
    
	return nil
}

func handleDevToArticle(article interface{}, user *types.User, integration *types.Integration) error {
    // get the slug of the artcile
    slug, ok := article.(map[string]interface{})["slug"].(string)
    if !ok {
        return errors.New("Couldn't decode slug from article")
    }

    // check if there is a fetch_link with the same identifier and fetch_source
    link, err := database.GetFetchLinkBySlug(integration.ID, slug)
    if err != nil {
        return err
    }

    if link != nil {
        // we already have this article
        // TODO: update existing post
        return nil
    }

    // check if there is a post with the same name
    title, ok := article.(map[string]interface{})["title"].(string)
    if !ok {
        return errors.New("Couldn't decode title from article")
    }

    post, err := database.GetPostsByTitle(user, title)
    if err != nil {
        return err
    }

    if post != nil {
        // we already have this article
        // create fetch_link
        new_link := types.FetchLink{
            PostID: post.Id,
            FetchSource: integration.ID,
            Identifier: slug,
        }
        err = database.CreateFetchLink(new_link)
        if err != nil {
            return err
        }
        return nil
    }

    // create a new post
    newPost := types.Post{
        Title: title,
        Content: article.(map[string]interface{})["body_markdown"].(string),
        AuthorID: user.ID,
        Archived: false,
        Slug: slug,
        Category: "devlog",
        Tags: []string{},
        Description: article.(map[string]interface{})["description"].(string),
    }
    post_id, err := database.CreatePost(newPost)
    if err != nil {
        return err
    }

    // create fetch_link
    new_link := types.FetchLink{
        PostID: post_id,
        FetchSource: integration.ID,
        Identifier: slug,
    }
    err = database.CreateFetchLink(new_link)
    if err != nil {
        return err
    }

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
