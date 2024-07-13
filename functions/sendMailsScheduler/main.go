package handler

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/functions"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/appwrite/sdk-for-go/users"
	"github.com/open-runtimes/types-for-go/v4"
)

func Main(Context *types.Context) types.ResponseOutput {
	appwriteClient := client.NewClient()
	appwriteClient.SetEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT"))
	appwriteClient.SetProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID"))
	appwriteClient.SetKey(Context.Req.Headers["x-appwrite-key"])

	appwriteUsers := users.NewUsers(appwriteClient)
	appwriteFunctions := functions.NewFunctions(appwriteClient)

	hour, minute, _ := time.Now().Clock()
	weekday := int(time.Now().Weekday())
	_, _, monthday := time.Now().Date()

	format := "am"
	if hour >= 12 {
		hour -= 12
		format = "pm"
	}

	today := time.Now()
	endOfMonthDay := today.AddDate(0, 1, -today.Day()).Day()
	beforeEndOfMonthDay := today.AddDate(0, 1, -today.Day()-1).Day()

	monthlyVerbose := ""

	if monthday == 7 {
		monthlyVerbose = "day1"
	} else if monthday == 14 {
		monthlyVerbose = "day7"
	} else if monthday == 21 {
		monthlyVerbose = "day14"
	} else if monthday == endOfMonthDay {
		monthlyVerbose = "dayLast"
	} else if monthday == beforeEndOfMonthDay {
		monthlyVerbose = "dayBeforeLast"
	}

	currentDateStrings := []string{
		strconv.Itoa(hour) + "T" + strconv.Itoa(minute) + "T" + format + "TdailyT",
		strconv.Itoa(hour) + "T" + strconv.Itoa(minute) + "T" + format + "TweeklyT" + strconv.Itoa(weekday),
	}

	if monthlyVerbose != "" {
		currentDateStrings = append(currentDateStrings, strconv.Itoa(hour)+"T"+strconv.Itoa(minute)+"T"+format+"TmonthlyT"+monthlyVerbose)
	}

	Context.Log(currentDateStrings)

	cursor := "INIT"
	for ok := true; ok; ok = (cursor != "") {
		Context.Log("Page iteration")

		orQueries := []string{}
		for _, currentDate := range currentDateStrings {
			orQueries = append(orQueries, query.Contains("labels", currentDate))
		}

		queries := []interface{}{
			query.Limit(50),
			query.Or(orQueries),
		}

		if cursor != "INIT" {
			queries = append(queries, query.CursorAfter(cursor))
		}

		listResponse, listErr := appwriteUsers.List(users.WithListQueries(queries))
		if listErr != nil {
			Context.Error(listErr)
			return Context.Res.Text("Error", 500, nil)
		}

		var wg sync.WaitGroup
		errCh := make(chan error, len(listResponse.Users))

		for _, user := range listResponse.Users {
			userStruct := user.(map[string]interface{})
			userLabels := userStruct["labels"].([]interface{})
			userLabel := userLabels[0].(string)

			hasLabel := false

			for _, currentDate := range currentDateStrings {
				if userLabel == currentDate {
					hasLabel = true
					break
				}
			}

			if !hasLabel {
				continue
			}

			wg.Add(1)
			go func(u interface{}) {
				defer wg.Done()

				userStruct := u.(map[string]interface{})
				userId := userStruct["$id"].(string)
				userEmail := userStruct["email"].(string)

				Context.Log("Sending mail to " + userId + ": " + userEmail)

				_, err := appwriteFunctions.CreateExecution(
					"sendMails",
					functions.WithCreateExecutionAsync(true),
					functions.WithCreateExecutionMethod("POST"),
					functions.WithCreateExecutionBody(userId),
				)
				if err != nil {
					errCh <- err
				}
			}(user)
		}

		wg.Wait()
		close(errCh)

		for err := range errCh {
			Context.Error(err)
			return Context.Res.Text("Error", 500, nil)
		}

		if len(listResponse.Users) > 0 {
			lastDocument := listResponse.Users[len(listResponse.Users)-1].(map[string]interface{})

			lastDocumentId := lastDocument["$id"].(string)
			cursor = lastDocumentId
		} else {
			cursor = ""
		}
	}

	Context.Log("Done")

	return Context.Res.Text("OK", 200, nil)
}
