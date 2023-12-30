package utils

import (
	"blog-server/types"
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
)

func LogError(message string, err error, status int, writer http.ResponseWriter) {
	log.Error(message, "err", err)
	http.Error(writer, message, status)
}

func Unauthorized(writer http.ResponseWriter) {
    http.Error(writer, "Unauthorized access", http.StatusUnauthorized);
}

func ResponseJSON(data interface{}, writer http.ResponseWriter) {
	encoded, err := json.Marshal(data)
	if err != nil {
		LogError("Error encoding to JSON", err, http.StatusInternalServerError, writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(encoded)
}

func GetChildrenCategories(categories []types.Category, parent string) []types.Category {
	var cats []types.Category

	for i := 0; i < len(categories); i++ {
		if categories[i].Parent == parent {
			cats = append(cats, categories[i])
			// add all the children as well
			var children = GetChildrenCategories(categories, categories[i].Name)
			for j := 0; j < len(children); j++ {
				cats = append(cats, children[j])
			}
		}
	}

	return cats
}

func GetUser(r *http.Request) *types.User {
    user, ok := r.Context().Value("user").(*types.User);
    if ok && user != nil {
        return user
    } else {
        return nil
    }
}
