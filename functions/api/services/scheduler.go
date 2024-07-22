package services

import (
	"strconv"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/users"
	openruntimes "github.com/open-runtimes/types-for-go/v4"
)

type UpdateSchedulerIntervalBody struct {
	Minutes          int    `json:"minutes"`
	Hours            int    `json:"hours"`
	Format           string `json:"format"`
	Frequency        string `json:"frequency"`
	FrequencyDetails string `json:"frequencyDetails"`
}

func UpdateSchedulerInterval(Context openruntimes.Context, appwriteClient client.Client) openruntimes.Response {
	var body UpdateSchedulerIntervalBody
	err := Context.Req.BodyJson(&body)
	if err != nil {
		return Context.Res.Text("Invalid body.", Context.Res.WithStatusCode(400))
	}

	// Validators
	if body.Format != "am" && body.Format != "pm" {
		return Context.Res.Text("Format must be 'am' or 'pm'", Context.Res.WithStatusCode(400))
	}

	if body.Frequency != "daily" && body.Frequency != "weekly" && body.Frequency != "monthly" {
		return Context.Res.Text("Frequency must be 'daily' or 'weekly' or 'monthly'", Context.Res.WithStatusCode(400))
	}

	if body.Minutes < 0 || body.Minutes > 59 {
		return Context.Res.Text("Minutes must be between 0 and 59", Context.Res.WithStatusCode(400))
	}

	if body.Hours < 0 || body.Hours > 11 {
		return Context.Res.Text("Hours must be between 0 and 11", Context.Res.WithStatusCode(400))
	}

	if body.Frequency == "weekly" {
		dayOfWeek, err := strconv.Atoi(body.FrequencyDetails)
		if err != nil || dayOfWeek < 0 || dayOfWeek > 6 {
			return Context.Res.Text("When frequency is 'weekly', frequencyDetails must be a number between 0 and 6", Context.Res.WithStatusCode(400))
		}
	}

	if body.Frequency == "monthly" {
		if body.FrequencyDetails != "day1" &&
			body.FrequencyDetails != "day7" &&
			body.FrequencyDetails != "day14" &&
			body.FrequencyDetails != "dayLast" &&
			body.FrequencyDetails != "dayBeforeLast" {
			return Context.Res.Text("When frequency is 'monthly', frequencyDetails must be 'day1' or 'day7' or 'day14' or 'dayLast' or 'dayBeforeLast'", Context.Res.WithStatusCode(400))
		}
	}

	if body.Frequency == "weekly" {
	}

	// Ensure it's user-executed
	userId, userIdOk := Context.Req.Headers["x-appwrite-user-id"]
	if !userIdOk || userId == "" {
		return Context.Res.Text("Unauthorized", Context.Res.WithStatusCode(401))
	}

	// Action
	label := strconv.Itoa(body.Hours) + "T" + strconv.Itoa(body.Minutes) + "T" + body.Format + "T" + body.Frequency + "T" + body.FrequencyDetails

	appwriteUsers := users.NewUsers(appwriteClient)
	_, err = appwriteUsers.UpdateLabels(userId, []interface{}{
		label,
	})

	if err != nil {
		return Context.Res.Text("Error updating user labels", Context.Res.WithStatusCode(500))
	}

	return Context.Res.Text("OK")
}
