package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/permission"
	"github.com/appwrite/sdk-for-go/role"
	"github.com/open-runtimes/types-for-go/v4"
)

type CreateCalendarBody struct {
	Url   string `json:"url"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func CreateCalendar(Context *types.Context, appwriteClient client.Client) types.ResponseOutput {
	var body CreateCalendarBody
	err := json.Unmarshal(Context.Req.BodyBinary(), &body)
	if err != nil {
		return Context.Res.Text("Invalid body.", 400, nil)
	}

	// Transformers
	if strings.HasPrefix(body.Url, "webcal://") {
		body.Url = "https://" + body.Url[len("webcal://"):]
	}

	// Validators
	if body.Color != "pink" && body.Color != "orange" && body.Color != "blue" && body.Color != "yellow" && body.Color != "purple" && body.Color != "mint" {
		return Context.Res.Text("Color must be 'pink' or 'orange' or 'blue' or 'yellow' or 'purple'", 400, nil)
	}

	_, err = http.Get(body.Url)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Calendar URL is not valid", 400, nil)
	}

	// Ensure it's user-executed
	userId, userIdOk := Context.Req.Headers["x-appwrite-user-id"]
	if !userIdOk || userId == "" {
		return Context.Res.Text("Unauthorized", 401, nil)
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
		return Context.Res.Text(err.Error(), 400, nil)
	}

	return Context.Res.Text("OK", 200, nil)
}
