package services

import (
	"github.com/appwrite/sdk-for-go/appwrite"
	c "github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

type DeleteSubscriptionBody struct {
	Email string `json:"email"`
}

type CustomUser struct {
	*models.User
	Prefs map[string]interface{} `json:"prefs"`
}

type CustomUserList struct {
	*models.UserList
	Users []CustomUser `json:"users"`
}

// TODO: Security: Use encrypted email values
func DeleteSubscription(Context openruntimes.Context, client c.Client) openruntimes.Response {
	var body DeleteSubscriptionBody
	err := Context.Req.BodyJson(&body)
	if err != nil {
		return Context.Res.Text("Invalid body.", Context.Res.WithStatusCode(400))
	}

	// Action
	users := appwrite.NewUsers(client)

	response, responseErr := users.List(users.WithListQueries([]string{
		query.Equal("email", body.Email),
		query.Limit(1),
	}))

	if responseErr != nil {
		return Context.Res.Text(responseErr.Error(), Context.Res.WithStatusCode(400))
	}

	// Return OK here to prevent account existance exposure
	if len(response.Users) == 0 {
		return Context.Res.Text("OK")
	}

	var typedUsers CustomUserList
	err = response.Decode(&typedUsers)
	if err != nil {
		return Context.Res.Text(err.Error(), Context.Res.WithStatusCode(400))
	}

	currentUser := typedUsers.Users[0]
	currentPrefs := currentUser.Prefs
	currentPrefs["unsubscribed"] = true

	_, err = users.UpdatePrefs(currentUser.Id, currentPrefs)
	if err != nil {
		return Context.Res.Text("Could not mark as unsubscribed", Context.Res.WithStatusCode(400))
	}

	return Context.Res.Text("OK")
}
