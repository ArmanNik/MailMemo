package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/apognu/gocal"
	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/permission"
	"github.com/appwrite/sdk-for-go/role"
)

type CreateCalendarBody struct {
	Url   string `json:"url"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func CreateCalendar(context openRuntimes.Context, client appwrite.Client) openRuntimes.ResponseOutput {
	var body CreateCalendarBody
	err := context.Req.BodyJson(&body)
	if err != nil {
		return context.Res.Text("Invalid body.", context.Res.WithStatusCode(400))
	}

	// Transformers
	if strings.HasPrefix(body.Url, "webcal://") {
		body.Url = "https://" + body.Url[len("webcal://"):]
	}

	// Validators
	if body.Color != "pink" && body.Color != "orange" && body.Color != "blue" && body.Color != "yellow" && body.Color != "purple" && body.Color != "mint" {
		return context.Res.Text("Color must be 'pink' or 'orange' or 'blue' or 'yellow' or 'purple'", context.Res.WithStatusCode(400))
	}

	calResp, err := http.Get(body.Url)
	if err != nil {
		context.Error(err)
		return context.Res.Text("Calendar URL is not valid", context.Res.WithStatusCode(400))
	}

	defer calResp.Body.Close()

	start := time.Now()
	end := time.Now().AddDate(10, 0, 0)

	c := gocal.NewParser(calResp.Body)
	c.Start, c.End = &start, &end
	parseErr := c.Parse()
	if parseErr != nil {
		context.Error(parseErr)
		return context.Res.Text("URL is not a valid calendar", context.Res.WithStatusCode(400))
	}

	if len(c.Events) == 0 {
		context.Error("No events found")
		return context.Res.Text("URL is not a valid calendar", context.Res.WithStatusCode(400))
	}

	// Ensure it's user-executed
	userId, userIdOk := context.Req.Headers["x-appwrite-user-id"]
	if !userIdOk || userId == "" {
		return context.Res.Text("Unauthorized", context.Res.WithStatusCode(401))
	}

	// Action
	databases := appwrite.NewDatabases(client)
	_, err = databases.CreateDocument("main", "calendars", id.Unique(), map[string]interface{}{
		"url":    body.Url,
		"name":   body.Name,
		"color":  body.Color,
		"userId": userId,
	}, databases.WithCreateDocumentPermissions([]string{
		permission.Read(role.User(userId, "")),
	}))

	if err != nil {
		return context.Res.Text(err.Error(), context.Res.WithStatusCode(400))
	}

	return context.Res.Text("OK")
}
