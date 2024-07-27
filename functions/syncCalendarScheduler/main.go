package handler

import (
	"os"
	"sync"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// START-OF-COPY-PASTE

// Appwrite User types
type AppwriteUserPrefs struct {
	Timezone     string `json:"timezone"`
	Period       string `json:"period"`
	Unsubscribed bool   `json:"unsubscribed"`
	FirstCal     bool   `json:"firstCal"`
	Onboarded    bool   `json:"onboarded"`
}

type AppwriteUser struct {
	*models.User
	Prefs AppwriteUserPrefs `json:"prefs"`
}

type AppwriteUserList struct {
	*models.UserList
	Users []AppwriteUser `json:"users"`
}

// Appwrite Calendar types
type AppwriteCalendarList struct {
	*models.DocumentList
	Documents []AppwriteCalendar `json:"documents"`
}

type AppwriteCalendar struct {
	*models.Document
	Name   string `json:"name"`
	Color  string `json:"color"`
	Url    string `json:"url"`
	UserId string `json:"userId"`
}

// Appwrite Event types

type AppwriteEventList struct {
	*models.DocumentList
	Documents []AppwriteEvent `json:"documents"`
}

type AppwriteEvent struct {
	*models.Document
	Name       string `json:"name"`
	Uid        string `json:"uid"`
	CalendarId string `json:"calendarId"`
	StartAt    string `json:"startAt"`
	EndAt      string `json:"endAt"`
	ModifiedAt string `json:"modifiedAt"`
}

// END-OF-COPY-PASTE

func Main(Context openruntimes.Context) openruntimes.Response {
	if Context.Req.Method != "POST" {
		return Context.Res.Text("Not Found", Context.Res.WithStatusCode(404))
	}

	client := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID")),
		appwrite.WithKey(Context.Req.Headers["x-appwrite-key"]),
	)

	databases := appwrite.NewDatabases(client)
	functions := appwrite.NewFunctions(client)

	userId := Context.Req.Headers["x-appwrite-user-id"]
	if userId == "" {
		userId = Context.Req.BodyText()
	}

	cursor := "INIT"
	for ok := true; ok; ok = (cursor != "") {
		Context.Log("Page iteration")

		queries := []string{
			query.Limit(50),
			query.Select([]interface{}{"$id"}),
		}

		if cursor != "INIT" {
			queries = append(queries, query.CursorAfter(cursor))
		}

		if userId != "" {
			queries = append(queries, query.Equal("userId", userId))
		}

		listResponse, listErr := databases.ListDocuments("main", "calendars", databases.WithListDocumentsQueries(queries))
		if listErr != nil {
			Context.Error(listErr)
			return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
		}

		var wg sync.WaitGroup
		wg.Add(len(listResponse.Documents))
		errCh := make(chan error, len(listResponse.Documents))

		for _, document := range listResponse.Documents {
			id := document.Id

			go func(id string) {
				defer wg.Done()

				Context.Log("Executing sync for calendar " + id)

				_, err := functions.CreateExecution(
					"syncCalendar",
					functions.WithCreateExecutionAsync(true),
					functions.WithCreateExecutionMethod("POST"),
					functions.WithCreateExecutionBody(id),
				)
				if err != nil {
					errCh <- err
				}
			}(id)
		}

		wg.Wait()
		close(errCh)

		for err := range errCh {
			Context.Error(err)
			return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
		}

		if len(listResponse.Documents) > 0 {
			lastDocument := listResponse.Documents[len(listResponse.Documents)-1]
			lastDocumentId := lastDocument.Id
			cursor = lastDocumentId
		} else {
			cursor = ""
		}
	}
	Context.Log("Done")

	return Context.Res.Text("OK")
}
