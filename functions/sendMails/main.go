package handler

// TODO: Respect period from prefs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Boostport/mjml-go"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/messaging"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/appwrite/sdk-for-go/users"
	"github.com/open-runtimes/types-for-go/v4"
)

type MailElement struct {
	Name          string
	CalendarName  string
	CalendarColor string
	Time          string
}

type MailDay struct {
	DaysRelative int
	Events       []MailElement
}

type MailData struct {
	Title     string
	Today     []MailElement
	Upcomming map[string]MailDay
}

func Main(Context *types.Context) types.ResponseOutput {
	if Context.Req.Method != "POST" {
		return Context.Res.Text("Not Found", 404, nil)
	}

	userId := Context.Req.Headers["x-appwrite-user-id"]

	if userId == "" {
		userId = Context.Req.BodyText()
	}

	appwriteClient := client.NewClient()
	appwriteClient.SetEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT"))
	appwriteClient.SetProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID"))
	appwriteClient.SetKey(Context.Req.Headers["x-appwrite-key"])

	appwriteMessaging := messaging.NewMessaging(appwriteClient)
	appwriteUsers := users.NewUsers(appwriteClient)
	appwriteDatabases := databases.NewDatabases(appwriteClient)

	userStruct, userStructErr := appwriteUsers.Get(userId)
	if userStructErr != nil {
		Context.Error(userStructErr)
		return Context.Res.Text("Error", 500, nil)
	}

	listCalendars, listCalendarsErr := appwriteDatabases.ListDocuments("main", "calendars", databases.WithListDocumentsQueries([]interface{}{
		query.Limit(100),
		query.Equal("userId", userId),
		query.Select([]interface{}{
			"name",
			"color",
			"$id",
		}),
	}))
	if listCalendarsErr != nil {
		Context.Error(listCalendarsErr)
		return Context.Res.Text("Error", 500, nil)
	}

	calendarIds := []interface{}{}
	for _, calendar := range listCalendars.Documents {
		calendarDocument := calendar.(map[string]interface{})
		calendarIds = append(calendarIds, calendarDocument["$id"].(string))
	}

	userPrefs := userStruct.Prefs.(map[string]interface{})
	timezoneString := userPrefs["timezone"].(string)
	userTimezone, timezoneErr := time.LoadLocation(timezoneString)

	if timezoneErr != nil {
		Context.Error(timezoneErr)
		return Context.Res.Text("Error", 500, nil)
	}

	currentTime := time.Now().In(userTimezone)
	currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	listEvents, listEventsErr := appwriteDatabases.ListDocuments("main", "events", databases.WithListDocumentsQueries([]interface{}{
		query.Limit(10000),
		query.OrderAsc("startAt"),
		query.Equal("calendarId", calendarIds),
		query.GreaterThanEqual("startAt", currentTime.Format(time.RFC3339)),
	}))
	if listEventsErr != nil {
		Context.Error(listEventsErr)
		return Context.Res.Text("Error", 500, nil)
	}

	mailData := MailData{
		Title:     time.Now().In(userTimezone).Format("Monday, January 2, 2006"),
		Today:     []MailElement{},
		Upcomming: map[string]MailDay{},
	}

	totalEventsInWeek := 0
	totalEventsInMonth := 0

	todayDate := time.Now().In(userTimezone).Format(time.DateOnly)

	for _, event := range listEvents.Documents {
		eventDocument := event.(map[string]interface{})
		eventStartAt, _ := time.Parse(time.RFC3339, eventDocument["startAt"].(string))
		eventEndAt, _ := time.Parse(time.RFC3339, eventDocument["endAt"].(string))
		eventDayKey := eventStartAt.Format(time.DateOnly)

		name := eventDocument["name"].(string)
		calendarId := eventDocument["calendarId"].(string)
		calendarName := ""
		calendarColor := ""
		for _, calendar := range listCalendars.Documents {
			calendarDocument := calendar.(map[string]interface{})
			if calendarDocument["$id"].(string) == calendarId {
				calendarName = calendarDocument["name"].(string)
				calendarColor = calendarDocument["color"].(string)
				break
			}
		}

		formattedTime := eventStartAt.Format("15:04") + " - " + eventEndAt.Format("15:04")
		if formattedTime == "00:00 - 23:59" {
			formattedTime = "All day"
		}

		element := MailElement{
			Name:          name,
			CalendarName:  calendarName,
			CalendarColor: calendarColor,
			Time:          formattedTime,
		}

		daysRelative := int(eventStartAt.Sub(currentTime).Hours() / 24)

		if eventDayKey == todayDate {
			mailData.Today = append(mailData.Today, element)
		} else {
			mailDay, ok := mailData.Upcomming[eventDayKey]

			if !ok {
				mailData.Upcomming[eventDayKey] = MailDay{
					DaysRelative: daysRelative,
					Events:       []MailElement{element},
				}
			} else {
				mailData.Upcomming[eventDayKey] = MailDay{
					DaysRelative: mailDay.DaysRelative,
					Events:       append(mailDay.Events, element),
				}
			}

			if daysRelative <= 7 {
				totalEventsInWeek++
				Context.Log("Adding " + element.Name)
				Context.Log(daysRelative)
			} else if daysRelative <= 30 {
				totalEventsInMonth++
			}
		}
	}

	todayReminders := len(mailData.Today)

	verboseWord := "reminders"
	if todayReminders == 1 {
		verboseWord = "reminder"
	}

	remindersVerbose := strconv.Itoa(todayReminders) + " " + verboseWord
	if todayReminders == 0 {
		remindersVerbose = "No reminders"
	}

	subject := remindersVerbose + " on 9.7.2024"
	previewText := strconv.Itoa(totalEventsInWeek) + " more in a week, and " + strconv.Itoa(totalEventsInMonth) + " more in a month... Today is:"

	input := `
			<mjml>
			<mj-head>
				<mj-title>` + subject + `</mj-title>
				<mj-preview>` + previewText + `</mj-preview>
			</mj-head>
			<mj-body>
			  <mj-wrapper border="1px solid #EDEDF0" padding="50px 30px" background-color="#FAFAFB">
				<mj-section>
				  <mj-column>
					<mj-text font-size="20px" color="#56565C" font-family="helvetica" padding-left="0px">
					  ` + mailData.Title + `
					</mj-text>
				  </mj-column>
				</mj-section>

				<mj-section>
				  <mj-column>
					<mj-text font-size="30px" color="#19191C" font-family="helvetica" padding-left="0px">
						Today
					</mj-text>
				  </mj-column>
				</mj-section>

				<mj-section>
					<mj-column>
						<mj-text font-size="12px" color="#2D2D31" font-family="helvetica" padding-top="0px">
							` + fmt.Sprintf("%v", mailData.Today) + `
						</mj-text>
					</mj-column>
				</mj-section>

				<mj-section>
				  <mj-column>
					<mj-text font-size="30px" color="#19191C" font-family="helvetica" padding-left="0px">
						Upcoming
					</mj-text>
				  </mj-column>
				</mj-section>

				<mj-section>
					<mj-column>
						<mj-text font-size="12px" color="#2D2D31" font-family="helvetica" padding-top="0px">
							` + fmt.Sprintf("%v", mailData.Upcomming) + `
						</mj-text>
					</mj-column>
				</mj-section>
			  </mj-wrapper>
			</mj-body>
		  </mjml>
			`

	output, err := mjml.ToHTML(
		context.Background(),
		input,
		mjml.WithValidationLevel(mjml.Strict),
		mjml.WithBeautify(false),
		mjml.WithMinify(true),
		mjml.WithKeepComments(false),
	)

	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		errorAsJson, _ := json.Marshal(mjmlError)
		Context.Error(string(errorAsJson[:]))
	}

	message, err := appwriteMessaging.CreateEmail(
		id.Unique(),
		subject,
		output,
		messaging.WithCreateEmailHtml(true),
		messaging.WithCreateEmailUsers([]interface{}{
			userStruct.Id,
		}),
	)

	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", 500, nil)
	}

	Context.Log("Message ID: " + message.Id)

	return Context.Res.Text("OK", 200, nil)

}
