package handler

import (
	"os"
	"strconv"
	"sync"
	"time"

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
	client := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID")),
		appwrite.WithKey(Context.Req.Headers["x-appwrite-key"]),
	)

	users := appwrite.NewUsers(client)
	functions := appwrite.NewFunctions(client)

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

		queries := []string{
			query.Limit(50),
			query.Or(orQueries),
		}

		if cursor != "INIT" {
			queries = append(queries, query.CursorAfter(cursor))
		}

		usersResponse, err := users.List(users.WithListQueries(queries))
		if err != nil {
			Context.Error(err)
			return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
		}

		var users AppwriteUserList
		err = usersResponse.Decode(&users)
		if err != nil {
			Context.Error(err)
			return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
		}

		var wg sync.WaitGroup
		errCh := make(chan error, len(users.Users))

		for _, user := range users.Users {
			userLabel := user.Labels[0]
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
			go func(u AppwriteUser) {
				defer wg.Done()

				userId := u.Id
				userEmail := u.Email

				Context.Log("Sending mail to " + userId + ": " + userEmail)

				_, err := functions.CreateExecution(
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
			return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
		}

		if len(users.Users) > 0 {
			lastDocument := users.Users[len(users.Users)-1]
			cursor = lastDocument.Id
		} else {
			cursor = ""
		}
	}

	Context.Log("Done")

	return Context.Res.Text("OK")
}
