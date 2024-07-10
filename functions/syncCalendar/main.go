package handler

import (
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/apognu/gocal"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/permission"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/appwrite/sdk-for-go/role"
	"github.com/open-runtimes/types-for-go/v4"
)

type EventMinimal struct {
	Uid          string
	Summary      string
	Description  string
	Start        *time.Time
	End          *time.Time
	LastModified *time.Time
}

func Main(Context *types.Context) types.ResponseOutput {
	if Context.Req.Method != "POST" {
		return Context.Res.Text("Not Found", 404, nil)
	}

	appwriteClient := client.NewClient()
	appwriteClient.SetEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT"))
	appwriteClient.SetProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID"))
	appwriteClient.SetKey(Context.Req.Headers["x-appwrite-key"])

	appwriteDatabases := databases.NewDatabases(appwriteClient)

	calendarId := Context.Req.BodyText()

	if calendarId == "" {
		Context.Log("Body missing")
		return Context.Res.Text("Error", 500, nil)
	}

	// TODO: Use getDocument
	calendarResponse, err := appwriteDatabases.ListDocuments("main", "calendars", databases.WithListDocumentsQueries([]interface{}{
		query.Equal("$id", calendarId),
	}))

	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", 500, nil)
	}

	if len(calendarResponse.Documents) == 0 {
		Context.Log(errors.New("Calendar not found"))
		return Context.Res.Text("Error", 500, nil)
	}

	calendarDocument := calendarResponse.Documents[0].(map[string]interface{})
	calendarUrl := calendarDocument["url"].(string)
	userId := calendarDocument["userId"].(string)

	calResp, err := http.Get(calendarUrl)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", 500, nil)
	}

	defer calResp.Body.Close()

	start, end := time.Now(), time.Now().Add(12*30*24*time.Hour)

	c := gocal.NewParser(calResp.Body)
	c.Start, c.End = &start, &end
	c.Parse()

	Context.Log("Processing")

	eventChunk := [100]EventMinimal{}

	i := 0
	for _, e := range c.Events {
		eventChunk[i] = EventMinimal{
			Uid:          e.Uid,
			Summary:      e.Summary,
			Description:  e.Description,
			Start:        e.Start,
			End:          e.End,
			LastModified: e.LastModified,
		}

		if i == 99 {
			err = processEventsChunk(Context, userId, calendarId, appwriteDatabases, eventChunk)
			if err != nil {
				Context.Error(err)
				return Context.Res.Text("Error", 500, nil)
			}

			eventChunk = [100]EventMinimal{}
			i = 0
			continue
		}

		i++
	}

	err = processEventsChunk(Context, userId, calendarId, appwriteDatabases, eventChunk)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", 500, nil)
	}

	eventChunk = [100]EventMinimal{}

	cursor := "INIT"
	for ok := true; ok; ok = (cursor != "") {
		Context.Error("Existing page iteration")

		var queries []interface{}

		if cursor == "INIT" {
			queries = []interface{}{
				query.Select([]interface{}{"$id", "uid"}),
				query.Limit(1000),
			}
		} else {
			queries = []interface{}{
				query.Select([]interface{}{"$id", "uid"}),
				query.Limit(1000),
				query.CursorAfter(cursor),
			}
		}

		listResponse, listErr := appwriteDatabases.ListDocuments("main", "events", databases.WithListDocumentsQueries(queries))
		if listErr != nil {
			Context.Error(listErr)
			return Context.Res.Text("Error", 500, nil)
		}

		for _, document := range listResponse.Documents {
			eventDocument := document.(map[string]interface{})
			id := eventDocument["$id"].(string)
			uid := eventDocument["uid"].(string)

			found := false
			for _, e := range c.Events {
				if e.Uid == uid {
					found = true
					break
				}
			}

			if found == false {
				Context.Log("Deleting " + uid)

				_, err := appwriteDatabases.DeleteDocument("main", "events", id)
				if err != nil {
					Context.Error(err)
					return Context.Res.Text("Error", 500, nil)
				}
			}
		}

		if len(listResponse.Documents) > 0 {
			lastDocument := listResponse.Documents[len(listResponse.Documents)-1].(map[string]interface{})
			lastDocumentId := lastDocument["$id"].(string)
			cursor = lastDocumentId
		} else {
			cursor = ""
		}
	}

	Context.Log("Finished")

	return Context.Res.Text("OK", 200, nil)
}

func processEventsChunk(Context *types.Context, userId string, calendarId string, appwriteDatabases *databases.Databases, events [100]EventMinimal) error {
	// TODO: Figure out how to remove deleted event
	eventIds := []interface{}{}

	for _, event := range events {
		if event.Uid == "" {
			continue
		}

		eventIds = append(eventIds, event.Uid)
	}

	eventsResponse, err := appwriteDatabases.ListDocuments("main", "events", databases.WithListDocumentsQueries([]interface{}{
		query.Equal("uid", eventIds),
		query.Limit(100),
		query.Equal("calendarId", calendarId),
		query.Select([]interface{}{"$id", "uid", "modifiedAt"}),
	}))

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 100)

	for _, event := range events {
		if event.Uid == "" {
			continue
		}

		var existingEventDocument map[string]interface{} = nil

		eventId := event.Uid
		for _, document := range eventsResponse.Documents {
			eventDocument := document.(map[string]interface{})
			if eventDocument["uid"].(string) == eventId {
				existingEventDocument = eventDocument
			}
		}

		if existingEventDocument == nil {
			Context.Log("Inserting " + event.Uid)

			wg.Add(1)
			go func(e EventMinimal) {
				defer wg.Done()

				_, err := appwriteDatabases.CreateDocument("main", "events", id.Unique(), map[string]interface{}{
					"calendarId":  calendarId,
					"uid":         e.Uid,
					"name":        e.Summary,
					"description": e.Description,
					"startAt":     e.Start.Format(time.RFC3339),
					"endAt":       e.End.Format(time.RFC3339),
					"modifiedAt":  e.LastModified.Format(time.RFC3339),
				}, databases.WithCreateDocumentPermissions([]interface{}{
					permission.Read(role.User(userId, "")),
				}))

				if err != nil {
					errCh <- err
				}
			}(event)
		} else {
			newLastModified := event.LastModified
			oldLastModified, err := time.Parse(time.RFC3339, existingEventDocument["modifiedAt"].(string))

			if err != nil {
				return err
			}

			if newLastModified.After(oldLastModified) {
				Context.Log("Updaring " + event.Uid)

				wg.Add(1)
				go func(e EventMinimal) {
					defer wg.Done()

					_, err := appwriteDatabases.UpdateDocument("main", "events", existingEventDocument["$id"].(string), databases.WithUpdateDocumentData(map[string]interface{}{
						"name":        e.Summary,
						"description": e.Description,
						"startAt":     e.Start.Format(time.RFC3339),
						"endAt":       e.End.Format(time.RFC3339),
						"modifiedAt":  e.LastModified.Format(time.RFC3339),
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
