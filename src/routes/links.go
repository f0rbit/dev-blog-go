package routes;

import (
	"blog-server/actions"
	"blog-server/database"
	"blog-server/types"
	"blog-server/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

func GetUserIntegrations(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }

    integrations, err := database.GetUserIntegrations(user.ID);
    if err != nil {
        utils.LogError("Error fetching integrations", err, http.StatusInternalServerError, w);
        return;
    }

    utils.ResponseJSON(integrations, w);
}

func UpsertIntegrations(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }

    // read the 'source' and 'data' from the body
    type IntegrationInput struct {
        Source string `json:"source"`
        Data string `json:"data"`
    }
    var input IntegrationInput;

    err := json.NewDecoder(r.Body).Decode(&input);

    if err != nil {
        utils.LogError("Error decoding new integration", err, http.StatusBadRequest, w);
        return;
    }

    // create a new integration
    // create the 'location' from the given source
    var location = utils.GetLocation(input.Source);
    var integration = types.Integration{
        UserID: user.ID,
        Source: input.Source,
        Data: input.Data,
        Location: location,
    }

    // verify that the user_id of the integration is the logged in user
    if integration.UserID != user.ID {
        utils.LogError("Invalid userID", errors.New("Create Integration userID doesn't match userID"), http.StatusBadRequest, w);
        return;
    }

    err = database.UpsertIntegration(integration);
    if err != nil {
        utils.LogError("Error upserting integration", err, http.StatusInternalServerError, w);
        return;
    }

    w.WriteHeader(http.StatusOK);
}

func FetchIntegration(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }

    source := mux.Vars(r)["source"];
    log.Info("Fetching integration", "source", source);
    
    switch source {
    case "devto":
        err := actions.SyncUserDevTo(user.ID);
        if err != nil {
            utils.LogError("Error syncing devto", err, http.StatusInternalServerError, w);
            return;
        }
        break;
    default:
        utils.LogError("Invalid source", errors.New("Invalid source"), http.StatusBadRequest, w);
        return;
    }

    w.WriteHeader(http.StatusOK);
}
