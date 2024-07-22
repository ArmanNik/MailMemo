package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/apognu/gocal"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/permission"
	"github.com/appwrite/sdk-for-go/role"
	openruntimes "github.com/open-runtimes/types-for-go/v4"
)

type CreateCalendarBody struct {
	Url   string `json:"url"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func CreateCalendar(Context openruntimes.Context, appwriteClient client.Client) openruntimes.Response {
	var body CreateCalendarBody
	err := Context.Req.BodyJson(&body)
	if err != nil {
		return Context.Res.Text("Invalid body.", Context.Res.WithStatusCode(400))
	}

	// Transformers
	if strings.HasPrefix(body.Url, "webcal://") {
		body.Url = "https://" + body.Url[len("webcal://"):]
	}

	// Validators
	if body.Color != "pink" && body.Color != "orange" && body.Color != "blue" && body.Color != "yellow" && body.Color != "purple" && body.Color != "mint" {
		return Context.Res.Text("Color must be 'pink' or 'orange' or 'blue' or 'yellow' or 'purple'", Context.Res.WithStatusCode(400))
	}

	calResp, err := http.Get(body.Url)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Calendar URL is not valid", Context.Res.WithStatusCode(400))
	}

	defer calResp.Body.Close()

	start := time.Now()
	end := time.Now().AddDate(10, 0, 0)

	c := gocal.NewParser(calResp.Body)
	c.Start, c.End = &start, &end
	parseErr := c.Parse()
	if parseErr != nil {
		Context.Error(parseErr)
		return Context.Res.Text("URL is not a valid calendar", Context.Res.WithStatusCode(400))
	}

	if len(c.Events) == 0 {
		Context.Error("No events found")
		return Context.Res.Text("URL is not a valid calendar", Context.Res.WithStatusCode(400))
	}

	// Ensure it's user-executed
	userId, userIdOk := Context.Req.Headers["x-appwrite-user-id"]
	if !userIdOk || userId == "" {
		return Context.Res.Text("Unauthorized", Context.Res.WithStatusCode(401))
	}

	// Action
	appwriteDatabase := databases.NewDatabases(appwriteClient)
	_, err = appwriteDatabase.CreateDocument("main", "calendars", id.Unique(), map[string]interface{}{
		"url":    body.Url,
		"name":   body.Name,
		"color":  body.Color,
		"userId": userId,
	}, databases.WithCreateDocumentPermissions([]interface{}{
		permission.Read(role.User(userId, "")),
	}))

	if err != nil {
		return Context.Res.Text(err.Error(), Context.Res.WithStatusCode(400))
	}

	return Context.Res.Text("OK")
}
