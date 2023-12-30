// tags.go
package routes

import (
	"blog-server/database"
	"blog-server/utils"
	"errors"
	"net/http"
	"strconv"
)

func AddPostTag(w http.ResponseWriter, r *http.Request) {
	params, err := getTagParams(r);
	if err != nil {
        utils.LogError("Error parsing params", err, http.StatusBadRequest, w);
		return;
	}

    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }
    // check if user owns the post
    post, err := database.FetchPost(user, database.ID, params.id);
    if err != nil {
        utils.LogError("Error fetching post", err, http.StatusInternalServerError, w);
        return;
    }
    if post.AuthorID != user.ID {
        utils.LogError("Unauthorized access", errors.New("userID doesn't match authorID in AddPostTag"), http.StatusUnauthorized, w);
        return;
    }

	err = database.CreateTag(params.id, params.tag);

	if err != nil {
        utils.LogError("Error inserting tag", err, http.StatusInternalServerError, w);
		return;
	}

	w.WriteHeader(http.StatusOK);
}

func DeletePostTag(w http.ResponseWriter, r *http.Request) {
	params, err := getTagParams(r)
	if err != nil {
        utils.LogError("Error parsing params", err, http.StatusBadRequest, w);
		return;
	}
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }
    // check if user owns the post
    post, err := database.FetchPost(user, database.ID, params.id);
    if err != nil {
        utils.LogError("Error fetching post", err, http.StatusInternalServerError, w);
        return;
    }
    if post.AuthorID != user.ID {
        utils.LogError("Unauthorized access", errors.New("userID doesn't match authorID in AddPostTag"), http.StatusUnauthorized, w);
        return;
    }

	err = database.DeleteTag(params.id, params.tag)

	if err != nil {
        utils.LogError("Error deleting tag", err, http.StatusInternalServerError, w);
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetTags(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }
	tags, err := database.GetTags(user)
	if err != nil {
        utils.LogError("Error fetching tags", err, http.StatusInternalServerError, w);
		return
	}

    utils.ResponseJSON(tags, w);
}

type TagParams struct {
	id  int
	tag string
}

func getTagParams(r *http.Request) (TagParams, error) {
	idstr := r.URL.Query().Get("id") // will be converted to an integer
	tag := r.URL.Query().Get("tag")
    params := TagParams{0, tag}
	if idstr != "" {
        id, err := strconv.Atoi(idstr)
		if err != nil {
			return params, err
		}
        params.id = id;
	} else {
        // id is a required param
		return params, errors.New("Invalid id")
	}

	if params.tag == "" || len(params.tag) <= 0 {
        // tag is a required param
		return params, errors.New("Invalid tag")
	}

	return params, nil
}
