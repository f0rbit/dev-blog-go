package routes

import (
	"blog-server/database"
	"blog-server/types"
	"blog-server/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

// GET /tokens
func GetUserTokens(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUser(r)

	if user == nil {
		utils.Unauthorized(w)
		return
	}

	tokens, err := database.GetTokens(user.ID)

	if err != nil {
		utils.LogError("Error getting tokens", err, http.StatusInternalServerError, w)
		return
	}

	utils.ResponseJSON(tokens, w)
}

// POST /token/new
func CreateToken(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }

    var newToken types.AccessKey;
    err := json.NewDecoder(r.Body).Decode(&newToken)
    if err != nil {
        utils.LogError("Error decoding new post", err, http.StatusBadRequest, w)
        return;
    }

    // verify that the user_id of the token is the logged in user
    if newToken.UserID != user.ID {
        utils.LogError("Invalid userID", errors.New("Create Token userID doesn't match userID"), http.StatusBadRequest, w);
        return;
    }

    id, err := database.CreateToken(newToken)

    if err != nil {
        utils.LogError("Error creating new token", err, http.StatusInternalServerError, w);
        return;
    }
    
    fetchedToken, err := database.GetToken(id);
    if err != nil {
        utils.LogError("Error fetching created token", err, http.StatusInternalServerError, w);
        return;
    }

    utils.ResponseJSON(fetchedToken, w)
}

func EditToken(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }

    var updateToken types.AccessKey;
    err := json.NewDecoder(r.Body).Decode(&updateToken)
    if err != nil {
        utils.LogError("Error decoding new post", err, http.StatusBadRequest, w)
        return;
    }

    // verify that the user_id of the token is the logged in user
    if updateToken.UserID != user.ID {
        utils.LogError("Invalid userID", errors.New("Create Token userID doesn't match userID"), http.StatusBadRequest, w);
        return;
    }

    log.Info("Updating token", "token", updateToken);

    err = database.UpdateToken(updateToken)
    if err != nil {
        utils.LogError("Error updating token", err, http.StatusInternalServerError, w);
        return;
    }

    w.WriteHeader(http.StatusOK);
}

func DeleteToken(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }

    tokenIDStr := mux.Vars(r)["id"]
    tokenID, err := strconv.Atoi(tokenIDStr);
    if err != nil {
        utils.LogError("Error parsing token ID", err, http.StatusBadRequest, w);
        return;
    }

    token, err := database.GetToken(tokenID);
    if err != nil {
        utils.LogError("Error fetching token", err, http.StatusInternalServerError, w);
        return;
    }

    if token.ID != tokenID {
        utils.LogError("Invalid tokenID", errors.New("Delete post authorID doesn't match user"), http.StatusBadRequest, w);
        return;
    }

    // delete the token
    err = database.DeleteToken(tokenID)
    if err != nil {
        utils.LogError("Error deleting token", err, http.StatusInternalServerError, w);
        return;
    }

    w.WriteHeader(http.StatusOK)
}
