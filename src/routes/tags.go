// tags.go
package routes

import (
    "blog-server/database"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
)

func AddPostTag(w http.ResponseWriter, r *http.Request) {
    id, tag, err := parseTagParams(r)
    if err != nil {
        log.Error("Error parsing params", "err", err)
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    err = database.CreateTag(id, tag);

    if err != nil {
        log.Warn("Error inserting tag", "err", err)
        http.Error(w, "Error creating tag", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func DeletePostTag(w http.ResponseWriter, r *http.Request) {
    id, tag, err := parseTagParams(r)
    if err != nil {
        log.Error("Error parsing params", "err", err)
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    err = database.DeleteTag(id, tag)

    if err != nil {
        log.Warn("Error deleting tag", "err", err)
        http.Error(w, "Error deleting tag", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func GetTags(w http.ResponseWriter, r *http.Request) {
    tags, err := database.GetTags();
    if err != nil {
        log.Error("Error fetching tags", "err", err);
        http.Error(w, "Internal Server Error", http.StatusInternalServerError);
        return;
    }
	// and then respond with json
	encoded, err := json.Marshal(tags)
	if err != nil {
		log.Error("Error marshalling categories to JSON", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return;
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func parseTagParams(r *http.Request) (int, string, error) {
	idstr := r.URL.Query().Get("id")
	tag := r.URL.Query().Get("tag")
	id := 0
	var err error
	if idstr != "" {
		id, err = strconv.Atoi(idstr)
		if err != nil {
			return 0, "", err
		}
	} else {
		return 0, "", errors.New("Invalid id")
	}

	if tag == "" || len(tag) <= 0 {
		return 0, "", errors.New("Invalid tag")
	}

	return id, tag, nil
}
