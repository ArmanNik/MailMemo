package services

import (
	"strconv"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/open-runtimes/types-for-go/v4"
)

type UpdateSchedulerIntervalBody struct {
	Minutes          int    `json:"minutes"`
	Hours            int    `json:"hours"`
	Format           string `json:"format"`
	Frequency        string `json:"frequency"`
	FrequencyDetails string `json:"frequencyDetails"`
}

func UpdateSchedulerInterval(context types.Context, client appwrite.Client) types.ResponseOutput {
	var body UpdateSchedulerIntervalBody
	err := context.Req.BodyJson(&body)
	if err != nil {
		return context.Res.Text("Invalid body.", context.Res.WithStatusCode(400))
	}

	// Validators
	if body.Format != "am" && body.Format != "pm" {
		return context.Res.Text("Format must be 'am' or 'pm'", context.Res.WithStatusCode(400))
	}

	if body.Frequency != "daily" && body.Frequency != "weekly" && body.Frequency != "monthly" {
		return context.Res.Text("Frequency must be 'daily' or 'weekly' or 'monthly'", context.Res.WithStatusCode(400))
	}

	if body.Minutes < 0 || body.Minutes > 59 {
		return context.Res.Text("Minutes must be between 0 and 59", context.Res.WithStatusCode(400))
	}

	if body.Hours < 0 || body.Hours > 11 {
		return context.Res.Text("Hours must be between 0 and 11", context.Res.WithStatusCode(400))
	}

	if body.Frequency == "weekly" {
		dayOfWeek, err := strconv.Atoi(body.FrequencyDetails)
		if err != nil || dayOfWeek < 0 || dayOfWeek > 6 {
			return context.Res.Text("When frequency is 'weekly', frequencyDetails must be a number between 0 and 6", context.Res.WithStatusCode(400))
		}
	}

	if body.Frequency == "monthly" {
		if body.FrequencyDetails != "day1" &&
			body.FrequencyDetails != "day7" &&
			body.FrequencyDetails != "day14" &&
			body.FrequencyDetails != "dayLast" &&
			body.FrequencyDetails != "dayBeforeLast" {
			return context.Res.Text("When frequency is 'monthly', frequencyDetails must be 'day1' or 'day7' or 'day14' or 'dayLast' or 'dayBeforeLast'", context.Res.WithStatusCode(400))
		}
	}

	if body.Frequency == "weekly" {
	}

	// Ensure it's user-executed
	userId, userIdOk := context.Req.Headers["x-appwrite-user-id"]
	if !userIdOk || userId == "" {
		return context.Res.Text("Unauthorized", context.Res.WithStatusCode(401))
	}

	// Action
	label := strconv.Itoa(body.Hours) + "T" + strconv.Itoa(body.Minutes) + "T" + body.Format + "T" + body.Frequency + "T" + body.FrequencyDetails

	users := appwrite.NewUsers(client)
	_, err = users.UpdateLabels(userId, []string{
		label,
	})

	if err != nil {
		return context.Res.Text("Error updating user labels", context.Res.WithStatusCode(500))
	}

	return context.Res.Text("OK")
}
