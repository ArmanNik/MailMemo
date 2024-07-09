package handler

import (
	"os"
	"sync"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/functions"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/open-runtimes/types-for-go/v4"
)

type CalendarDocument struct {
	*models.Document
	// Add more if needed
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
	appwriteFunctions := functions.NewFunctions(appwriteClient)

	cursor := "INIT"
	for ok := true; ok; ok = (cursor == "") {
		Context.Error("Page iteration")

		var queries []interface{}

		if cursor == "INIT" {
			queries = []interface{}{
				query.Limit(50),
			}
		} else {
			queries = []interface{}{
				query.Limit(50),
				query.CursorAfter(cursor),
			}
		}

		listResponse, listErr := appwriteDatabases.ListDocuments("main", "calendars", databases.WithListDocumentsQueries(queries))
		if listErr != nil {
			Context.Error(listErr)
			return Context.Res.Text("Error", 500, nil)
		}

		var wg sync.WaitGroup
		wg.Add(len(listResponse.Documents))

		for _, document := range listResponse.Documents {
			calendarDocument := document.(map[string]interface{})
			id := calendarDocument["$id"].(string)

			go func() {
				defer wg.Done()

				appwriteFunctions.CreateExecution(
					"syncCalendar",
					functions.WithCreateExecutionAsync(true),
					functions.WithCreateExecutionMethod("POST"),
					functions.WithCreateExecutionBody(id),
				)
			}()
		}

		wg.Wait()

		if len(listResponse.Documents) > 0 {
			lastDocument := listResponse.Documents[len(listResponse.Documents)-1].(map[string]interface{})
			lastDocumentId := lastDocument["$id"].(string)
			cursor = lastDocumentId
		} else {
			cursor = ""
		}
	}
	Context.Error("Done")

	return Context.Res.Text("OK", 200, nil)
}
