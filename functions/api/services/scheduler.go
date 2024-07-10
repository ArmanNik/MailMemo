package services

import (
	"encoding/json"
	"strconv"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/users"
	"github.com/open-runtimes/types-for-go/v4"
)

type UpdateSchedulerIntervalBody struct {
	Minutes int    `json:"minutes"`
	Hours   int    `json:"hours"`
	Format  string `json:"format"`
}

func UpdateSchedulerInterval(Context *types.Context, appwriteClient client.Client) types.ResponseOutput {
	var body UpdateSchedulerIntervalBody
	err := json.Unmarshal(Context.Req.BodyBinary(), &body)
	if err != nil {
		return Context.Res.Text("Invalid body.", 400, nil)
	}

	// Validators
	if body.Format != "am" && body.Format != "pm" {
		return Context.Res.Text("Format must be 'am' or 'pm'", 400, nil)
	}

	if body.Minutes < 0 || body.Minutes > 59 {
		return Context.Res.Text("Minutes must be between 0 and 59", 400, nil)
	}

	if body.Hours < 0 || body.Hours > 11 {
		return Context.Res.Text("Hours must be between 0 and 11", 400, nil)
	}

	// Ensure it's user-executed
	userId, userIdOk := Context.Req.Headers["x-appwrite-user-id"]
	if !userIdOk || userId == "" {
		return Context.Res.Text("Unauthorized", 401, nil)
	}

	// Action
	label := strconv.Itoa(body.Hours) + "T" + strconv.Itoa(body.Minutes) + "T" + body.Format

	appwriteUsers := users.NewUsers(appwriteClient)
	_, err = appwriteUsers.UpdateLabels(userId, []interface{}{
		label,
	})

	if err != nil {
		return Context.Res.Text("Error updating user labels", 500, nil)
	}

	return Context.Res.Text("OK", 200, nil)
}
