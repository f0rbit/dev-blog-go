package utils

import (
	"blog-server/types"
	"encoding/json"
	"html"
	"net/http"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/russross/blackfriday/v2"
)

func LogError(message string, err error, status int, writer http.ResponseWriter) {
	log.Error(message, "err", err)
	http.Error(writer, message, status)
}

func Unauthorized(writer http.ResponseWriter) {
	http.Error(writer, "Unauthorized access", http.StatusUnauthorized)
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
	user, ok := r.Context().Value("user").(*types.User)
	if ok && user != nil {
		return user
	} else {
		return nil
	}
}

func GetDescription(content string) string {
    // parses content markdown to html
	bytes := blackfriday.Run([]byte(content))
	text := string(bytes)

    // replaces all html tags
	r := regexp.MustCompile("<[^>]*>")
    // removes \n
	text = strings.ReplaceAll(text, "\n", " ")
	stripped := r.ReplaceAllString(text, "")
    // removes any html escaped characters
	stripped = html.UnescapeString(stripped)
    // limits the length
	description := firstNLinesOrChars(stripped, 3, 80)
    // adds "..."
	return description + "..."
}

// firstNLinesOrChars returns the first n lines or first numChars characters of the string, whichever is smaller.
func firstNLinesOrChars(s string, n, numChars int) string {
	var lineCount, charCount int
	for i, rune := range s {
		if rune == '\n' {
			lineCount++
		}
		if lineCount >= n || charCount >= numChars {
			return s[:i]
		}
		charCount++
	}
	// Return the entire string if it's shorter than the specified lengths
	return s
}
