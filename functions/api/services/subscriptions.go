package services

import (
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/appwrite/sdk-for-go/users"
	openruntimes "github.com/open-runtimes/types-for-go/v4"
)

type DeleteSubscriptionBody struct {
	Email string `json:"email"`
}

// TODO: Security: Use encrypted email values
func DeleteSubscription(Context *openruntimes.Context, appwriteClient client.Client) openruntimes.Response {
	var body DeleteSubscriptionBody
	err := Context.Req.BodyJson(&body)
	if err != nil {
		return Context.Res.Text("Invalid body.", Context.Res.WithStatusCode(400))
	}

	// Action
	appwriteUsers := users.NewUsers(appwriteClient)

	users, userErr := appwriteUsers.List(users.WithListQueries([]interface{}{
		query.Equal("email", body.Email),
		query.Limit(1),
	}))

	if userErr != nil {
		return Context.Res.Text(userErr.Error(), Context.Res.WithStatusCode(400))
	}

	// Return OK here to prevent account existance exposure
	if len(users.Users) == 0 {
		return Context.Res.Text("OK")
	}

	currentUser := users.Users[0].(map[string]interface{})
	currentPrefs := currentUser["prefs"].(map[string]interface{})

	currentPrefs["unsubscribed"] = true

	_, err = appwriteUsers.UpdatePrefs(currentUser["$id"].(string), currentPrefs)
	if err != nil {
		return Context.Res.Text("Could not mark as unsubscribed", Context.Res.WithStatusCode(400))
	}

	return Context.Res.Text("OK")
}
