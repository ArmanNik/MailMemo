package handler

import (
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/apognu/gocal"
	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/appwrite/sdk-for-go/permission"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/appwrite/sdk-for-go/role"
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

type EventMinimal struct {
	Uid          string
	Summary      string
	Start        *time.Time
	End          *time.Time
	LastModified *time.Time
}

func Main(Context openruntimes.Context) openruntimes.Response {
	if Context.Req.Method != "POST" {
		return Context.Res.Text("Not Found", Context.Res.WithStatusCode(404))
	}

	client := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID")),
		appwrite.WithKey(Context.Req.Headers["x-appwrite-key"]),
	)

	users := appwrite.NewUsers(client)
	databases := appwrite.NewDatabases(client)

	calendarId := Context.Req.BodyText()

	if calendarId == "" {
		Context.Log("Body missing")
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	Context.Log(calendarId)

	// TODO: Use getDocument
	calendarResponse, err := databases.ListDocuments("main", "calendars", databases.WithListDocumentsQueries([]string{
		query.Equal("$id", calendarId),
	}))

	if err != nil {
		Context.Error("Cannot get calendar:")
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	if len(calendarResponse.Documents) == 0 {
		Context.Log(errors.New("Calendar not found"))
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	var calendars AppwriteCalendarList
	err = calendarResponse.Decode(&calendars)

	calendarUrl := calendars.Documents[0].Url
	userId := calendars.Documents[0].UserId

	userResponse, err := users.Get(userId)
	if err != nil {
		Context.Error("Cannot get user:")
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	var user AppwriteUser
	err = userResponse.Decode(&user)
	if err != nil {
		Context.Error("Cannot decode user:")
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	timezoneString := user.Prefs.Timezone
	userTimezone, timezoneErr := time.LoadLocation(timezoneString)

	if timezoneErr != nil {
		Context.Error("Cannot load timezone:")
		Context.Error(timezoneErr)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	calResp, err := http.Get(calendarUrl)
	if err != nil {
		Context.Error("Cannot get calendar from URL:")
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	defer calResp.Body.Close()

	start := time.Now().In(userTimezone)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	end := time.Now().In(userTimezone)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	end = end.Add(12 * 30 * 24 * time.Hour)

	c := gocal.NewParser(calResp.Body)
	c.Start, c.End = &start, &end
	c.Parse()

	Context.Log("Processing")

	eventChunk := [50]EventMinimal{}

	i := 0
	for _, e := range c.Events {
		eventChunk[i] = EventMinimal{
			Uid:          e.Uid + e.Start.Format(time.RFC3339),
			Summary:      e.Summary,
			Start:        e.Start,
			End:          e.End,
			LastModified: e.LastModified,
		}

		if i == 49 {
			err = processEventsChunk(Context, userId, calendarId, databases, eventChunk)
			if err != nil {
				Context.Error("Cannot process chunk:")
				Context.Error(err)
				return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
			}

			eventChunk = [50]EventMinimal{}
			i = 0
			continue
		}

		i++
	}

	if len(eventChunk) > 0 {
		err = processEventsChunk(Context, userId, calendarId, databases, eventChunk)
		if err != nil {
			Context.Error("Cannot process final chunk:")
			Context.Error(err)
			return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
		}
	}

	eventChunk = [50]EventMinimal{}

	cursor := "INIT"
	for ok := true; ok; ok = (cursor != "") {
		Context.Log("Existing page iteration")

		var queries []string

		if cursor == "INIT" {
			queries = []string{
				query.Equal("calendarId", calendarId),
				query.Select([]interface{}{"$id", "uid"}),
				query.Limit(1000),
			}
		} else {
			queries = []string{
				query.Equal("calendarId", calendarId),
				query.Select([]interface{}{"$id", "uid"}),
				query.Limit(1000),
				query.CursorAfter(cursor),
			}
		}

		listResponse, err := databases.ListDocuments("main", "events", databases.WithListDocumentsQueries(queries))
		if err != nil {
			Context.Error("Cannot list documents during deletion:")
			Context.Error(err)
			return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
		}

		var events AppwriteEventList
		err = listResponse.Decode(&events)
		if err != nil {
			Context.Error("Cannot decode documents during deletion:")
			Context.Error(err)
			return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
		}

		for _, document := range events.Documents {
			id := document.Id
			uid := document.Uid

			found := false
			for _, e := range c.Events {
				eventUid := e.Uid + e.Start.Format(time.RFC3339)
				if eventUid == uid {
					found = true
					break
				}
			}

			if found == false {
				Context.Log("Deleting " + uid)

				_, err := databases.DeleteDocument("main", "events", id)
				if err != nil {
					Context.Error("Cannot delete event:")
					Context.Error(err)
					return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
				}
			}
		}

		if len(listResponse.Documents) > 0 {
			lastDocument := listResponse.Documents[len(listResponse.Documents)-1]
			lastDocumentId := lastDocument.Id
			cursor = lastDocumentId
		} else {
			cursor = ""
		}
	}

	Context.Log("Finished")

	return Context.Res.Text("OK")
}

func processEventsChunk(Context openruntimes.Context, userId string, calendarId string, databases *databases.Databases, events [50]EventMinimal) error {
	eventIds := []interface{}{}

	for _, event := range events {
		if event.Uid == "" {
			continue
		}

		eventIds = append(eventIds, event.Uid)
	}

	if len(eventIds) == 0 {
		return nil
	}

	eventsResponse, err := databases.ListDocuments("main", "events", databases.WithListDocumentsQueries([]string{
		query.Equal("uid", eventIds),
		query.Limit(50),
		query.Equal("calendarId", calendarId),
		query.Select([]interface{}{"$id", "uid", "modifiedAt"}),
	}))

	if err != nil {
		return err
	}

	var remoteEvents AppwriteEventList
	err = eventsResponse.Decode(&remoteEvents)

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 50)

	for _, event := range events {
		if event.Uid == "" {
			continue
		}

		var existingEventDocument AppwriteEvent = AppwriteEvent{}

		for _, document := range remoteEvents.Documents {
			if document.Uid == event.Uid {
				existingEventDocument = document
			}
		}

		if existingEventDocument == (AppwriteEvent{}) {
			Context.Log("Inserting " + event.Uid)

			wg.Add(1)
			go func(e EventMinimal) {
				defer wg.Done()

				_, err := databases.CreateDocument("main", "events", id.Unique(), map[string]interface{}{
					"calendarId": calendarId,
					"uid":        e.Uid,
					"name":       e.Summary,
					"startAt":    e.Start.Format(time.RFC3339),
					"endAt":      e.End.Format(time.RFC3339),
					"modifiedAt": e.LastModified.Format(time.RFC3339),
				}, databases.WithCreateDocumentPermissions([]string{
					permission.Read(role.User(userId, "")),
				}))

				if err != nil {
					errCh <- err
				}
			}(event)
		} else {
			newLastModified := event.LastModified
			oldLastModified, err := time.Parse(time.RFC3339, existingEventDocument.ModifiedAt)

			if err != nil {
				return err
			}

			if newLastModified.After(oldLastModified) {
				Context.Log("Updating " + event.Uid)

				wg.Add(1)
				go func(e EventMinimal) {
					defer wg.Done()

					_, err := databases.UpdateDocument("main", "events", existingEventDocument.Id, databases.WithUpdateDocumentData(map[string]interface{}{
						"name":       e.Summary,
						"startAt":    e.Start.Format(time.RFC3339),
						"endAt":      e.End.Format(time.RFC3339),
						"modifiedAt": e.LastModified.Format(time.RFC3339),
					}))

					if err != nil {
						errCh <- err
					}
				}(event)
			}
		}
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	return nil
}
