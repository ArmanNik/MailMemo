package services

import (
	"encoding/json"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/appwrite/sdk-for-go/users"
	"github.com/open-runtimes/types-for-go/v4"
)

type DeleteSubscriptionBody struct {
	Email string `json:"email"`
}

// TODO: Security: Use encrypted email values
func DeleteSubscription(Context *types.Context, appwriteClient client.Client) types.ResponseOutput {
	var body DeleteSubscriptionBody
	err := json.Unmarshal(Context.Req.BodyBinary(), &body)
	if err != nil {
		return Context.Res.Text("Invalid body.", 400, nil)
	}

	// Action
	appwriteUsers := users.NewUsers(appwriteClient)

	users, userErr := appwriteUsers.List(users.WithListQueries([]interface{}{
		query.Equal("email", body.Email),
		query.Limit(1),
	}))

	if userErr != nil {
		return Context.Res.Text(userErr.Error(), 400, nil)
	}

	// Return OK here to prevent account existance exposure
	if len(users.Users) == 0 {
		return Context.Res.Text("OK", 200, nil)
	}

	currentUser := users.Users[0].(map[string]interface{})
	currentPrefs := currentUser["prefs"].(map[string]interface{})

	currentPrefs["unsubscribed"] = true

	_, err = appwriteUsers.UpdatePrefs(currentUser["$id"].(string), currentPrefs)
	if err != nil {
		return Context.Res.Text("Could not mark as unsubscribed", 400, nil)
	}

	return Context.Res.Text("OK", 200, nil)
}
