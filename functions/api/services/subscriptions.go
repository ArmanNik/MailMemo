package services

import (
	"encoding/json"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/open-runtimes/types-for-go/v4"
)

type DeleteSubscriptionBody struct {
	Email string `json:"email"`
}

type UserPrefs struct {
	Subscribed bool `json:"subscribed"`
}

type UserWithPrefs struct {
	*appwrite.models.User
	Prefs UserPrefs
}

// TODO: Security: Use encrypted email values
func DeleteSubscription(ontext types.Context, client appwrite.Client) types.ResponseOutput {
	var body DeleteSubscriptionBody
	err := context.Req.BodyJson(&body)
	if err != nil {
		return context.Res.Text("Invalid body.", context.Res.WithStatusCode(400))
	}

	// Action
	users := appwrite.NewUsers(client)

	response, userErr := users.List[UserWithPrefs](users.WithListQueries([]string{
		query.Equal("email", body.Email),
		query.Limit(1),
	}))

	if userErr != nil {
		return context.Res.Text(userErr.Error(), context.Res.WithStatusCode(400))
	}

	// Return OK here to prevent account existance exposure
	if len(response.Users) == 0 {
		return context.Res.Text("OK")
	}

	currentUser := response.Users[0]
	currentPrefs := currentUser["prefs"]

	currentPrefs.Unsubscribed = true

	_, err = users.UpdatePrefs(currentUser.Id, currentPrefs)
	if err != nil {
		return context.Res.Text("Could not mark as unsubscribed", context.Res.WithStatusCode(400))
	}

	return context.Res.Text("OK")
}
