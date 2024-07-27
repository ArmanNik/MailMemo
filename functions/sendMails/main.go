package handler

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
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

// Function-specific types
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
	Title          string
	Today          []MailElement
	Upcomming      map[string]MailDay
	UpcommingOrder []string
}

func Main(Context openruntimes.Context) openruntimes.Response {
	if Context.Req.Method != "POST" {
		return Context.Res.Text("Not Found", Context.Res.WithStatusCode(404))
	}

	userId := Context.Req.Headers["x-appwrite-user-id"]

	if userId == "" {
		userId = Context.Req.BodyText()
	}

	client := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID")),
		appwrite.WithKey(Context.Req.Headers["x-appwrite-key"]),
	)

	functions := appwrite.NewFunctions(client)
	users := appwrite.NewUsers(client)
	databases := appwrite.NewDatabases(client)

	responseUser, err := users.Get(userId)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	var user AppwriteUser
	err = responseUser.Decode(&user)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	responseCalendars, err := databases.ListDocuments("main", "calendars", databases.WithListDocumentsQueries([]string{
		query.Limit(100),
		query.Equal("userId", userId),
		query.Select([]interface{}{
			"name",
			"color",
			"$id",
		}),
	}))
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	var listCalendars AppwriteCalendarList
	err = responseCalendars.Decode(&listCalendars)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	calendarIds := []interface{}{}
	for _, calendar := range listCalendars.Documents {
		calendarIds = append(calendarIds, calendar.Id)
	}

	userEmail := user.Email
	periodString := user.Prefs.Period
	userTimezone, timezoneErr := time.LoadLocation(user.Prefs.Timezone)

	if user.Prefs.Unsubscribed {
		Context.Log("User is unsubscribed")
		return Context.Res.Text("OK")
	}

	if timezoneErr != nil {
		Context.Error(timezoneErr)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	currentTime := time.Now().In(userTimezone)
	currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	endAt := time.Now().In(userTimezone)
	endAt = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	if periodString == "day" {
		endAt = endAt.Add(1 * 24 * time.Hour)
	} else if periodString == "week" {
		endAt = endAt.Add(7 * 24 * time.Hour)
	} else if periodString == "month" {
		endAt = endAt.Add(30 * 24 * time.Hour)
	} else if periodString == "months" {
		endAt = endAt.Add(3 * 30 * 24 * time.Hour)
	} else if periodString == "year" {
		endAt = endAt.Add(365 * 24 * time.Hour)
	}

	responseEvents, err := databases.ListDocuments("main", "events", databases.WithListDocumentsQueries([]string{
		query.Limit(10000),
		query.OrderAsc("startAt"),
		query.Equal("calendarId", calendarIds),
		query.GreaterThanEqual("startAt", currentTime.Format(time.RFC3339)),
		query.LessThan("endAt", endAt.Format(time.RFC3339)),
	}))
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	var events AppwriteEventList
	err = responseEvents.Decode(&events)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	mailData := MailData{
		Title:          time.Now().In(userTimezone).Format("Monday, January 2, 2006"),
		Today:          []MailElement{},
		Upcomming:      map[string]MailDay{},
		UpcommingOrder: []string{},
	}

	totalEventsInWeek := 0
	totalEventsInMonth := 0
	totalUpcoming := 0

	todayDate := time.Now().In(userTimezone).Format(time.DateOnly)

	for _, event := range events.Documents {
		eventStartAt, _ := time.Parse(time.RFC3339, event.StartAt)
		eventEndAt, _ := time.Parse(time.RFC3339, event.EndAt)
		eventDayKey := eventStartAt.In(userTimezone).Format(time.DateOnly)

		name := event.Name
		calendarId := event.CalendarId
		calendarName := ""
		calendarColor := ""
		for _, calendar := range listCalendars.Documents {
			if calendar.Id == calendarId {
				calendarName = calendar.Name
				calendarColor = calendar.Color
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
			totalUpcoming++

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

			if !arrayContains(mailData.UpcommingOrder, eventDayKey) {
				mailData.UpcommingOrder = append(mailData.UpcommingOrder, eventDayKey)
			}

			if daysRelative <= 7 {
				totalEventsInWeek++
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

	subject := remindersVerbose + " on " + currentTime.Format("02 Jan 2006")
	previewText := strconv.Itoa(totalEventsInWeek) + " more in a week, and " + strconv.Itoa(totalEventsInMonth) + " more in a month... Today is:"

	TodayHtml := `
		<mj-section padding="10px 0px 0px 0px" border-radius="16px 16px 16px 16px" css-class="small-wrapper">
			<mj-text font-size="16px" font-weight="400" color="#707076" padding="0px 0px 0px 0px" line-height="22px">
				No events to show
			</mj-text>
		</mj-section>
	`

	if len(mailData.Today) > 0 {
		TodayHtmlElements := []string{}

		i := 0
		for _, event := range mailData.Today {

			borderRadius := "16px 16px 16px 16px"

			if len(mailData.Today) == 2 {
				if i == 0 {
					borderRadius = "16px 16px 0px 0px"
				} else {
					borderRadius = "0px 0px 16px 16px"
				}
			} else if len(mailData.Today) > 2 {
				if i == 0 {
					borderRadius = "16px 16px 0px 0px"
				} else if i == len(mailData.Today)-1 {
					borderRadius = "0px 0px 16px 16px"
				} else {
					borderRadius = "0px 0px 0px 0px"
				}
			}

			TodayHtmlElements = append(TodayHtmlElements, `
				<mj-section background-color="#a8a8a9" padding="1px" border-radius="`+borderRadius+`" css-class="small-wrapper">
					<mj-column border-radius="`+borderRadius+`" background-color="#141416" padding="16px">
					<mj-text font-weight="300" font-size="12px" color="#C3C3C6" padding="0px 0px 12px 0px">
						`+event.Time+` &nbsp;&nbsp;<span style="color: #303031;">|</span>&nbsp;
						<div style="width: 8px; height: 8px; background-color: #`+getHex(event.CalendarColor)+`; border-radius: 50%; display: inline-block; vertical-align: middle; margin-top: -2px;"></div>
						&nbsp;`+event.CalendarName+`
					</mj-text>
					<mj-text font-size="18px" font-weight="400" color="#EDEDF0" padding="0px 0px 0px 0px" line-height="22px">
						`+event.Name+`
					</mj-text>
					</mj-column>
				</mj-section>
			`)

			i++
		}

		TodayHtml = strings.Join(TodayHtmlElements[:], `
			<mj-section padding="4px 0" css-class="small-wrapper"> </mj-section>
		`)
	}

	UpcomingHtml := `
		<mj-section padding="10px 0px 0px 0px" border-radius="16px 16px 16px 16px" css-class="small-wrapper">
			<mj-text font-size="16px" font-weight="400" color="#707076" padding="0px 0px 0px 0px" line-height="22px">
				No events to show
			</mj-text>
		</mj-section>
	`

	if len(mailData.UpcommingOrder) > 0 {
		UpcomingHtmlElements := []string{}

		for _, dayKey := range mailData.UpcommingOrder {
			day := mailData.Upcomming[dayKey]

			DayEventsHtmlElements := []string{}

			for _, event := range day.Events {
				DayEventsHtmlElements = append(DayEventsHtmlElements, `
					<mj-text font-weight="300" font-size="12px" color="#C3C3C6" padding="0px 0px 12px 0px">
						`+event.Time+` &nbsp;&nbsp;<span style="color: #303031;">|</span>&nbsp;
						<div style="width: 8px; height: 8px; background-color: #`+getHex(event.CalendarColor)+`; border-radius: 50%; display: inline-block; vertical-align: middle; margin-top: -2px;"></div>
						&nbsp;`+event.CalendarName+`
					</mj-text>
					<mj-text font-size="18px" font-weight="400" color="#EDEDF0" padding="0px 0px 0px 0px" line-height="22px">
						`+event.Name+`
					</mj-text>
				`)
			}

			DayEventsHtml := strings.Join(DayEventsHtmlElements[:], `
				<mj-divider border-width="16px" padding="0px 0px 0px 0px" border-color="#141416" />
				<mj-divider border-width="1px" padding="0px 0px 0px 0px" border-color="#222224" />
				<mj-divider border-width="16px" padding="0px 0px 0px 0px" border-color="#141416" />
			`)

			borderRadius := "16px 16px 16px 16px"

			previousDay, _ := time.Parse(time.DateOnly, dayKey)
			previousDay = previousDay.Add(-24 * time.Hour)
			previousDayKey := previousDay.Format(time.DateOnly)

			nextDay, _ := time.Parse(time.DateOnly, dayKey)
			nextDay = nextDay.Add(24 * time.Hour)
			nextDayKey := nextDay.Format(time.DateOnly)

			_, hasPreviousDay := mailData.Upcomming[previousDayKey]
			_, hasNextDay := mailData.Upcomming[nextDayKey]

			if hasPreviousDay && hasNextDay {
				borderRadius = "0px 0px 0px 0px"
			} else if hasPreviousDay {
				borderRadius = "0px 0px 16px 16px"
			} else if hasNextDay {
				borderRadius = "16px 16px 0px 0px"
			}

			dayVerbose := "days"
			if day.DaysRelative == 1 {
				dayVerbose = "day"
			}

			DividerHtml := `
				<mj-section padding="16px 0" css-class="small-wrapper"> </mj-section>
			`

			if hasNextDay {
				DividerHtml = `
					<mj-section padding="4px 0" css-class="small-wrapper"> </mj-section>
				`
			}

			chipColor := "#18181B"
			chipBgColor := "#FFFFFF"

			if day.DaysRelative > 7 {
				chipColor = "#707076"
				chipBgColor = "#222224"
			}

			UpcomingHtmlElements = append(UpcomingHtmlElements, `
				<mj-section background-color="#272729" padding="1px" border-radius="`+borderRadius+`" css-class="small-wrapper">
					<mj-column border-radius="`+borderRadius+`" background-color="#141416" padding="16px">
					
					<mj-text font-size="13px" font-weight="500" color="`+chipColor+`" padding="4px 0px 18px 0px">
						<span style="background-color: `+chipBgColor+`; border-radius:25px; padding: 5px 8px;">&nbsp;In `+strconv.Itoa(day.DaysRelative)+` `+dayVerbose+`&nbsp;</span>
					</mj-text>
					<mj-divider border-width="1px" padding="0px 0px 16px 0px" border-color="#222224" />
					
					`+DayEventsHtml+`
					</mj-column>
				</mj-section>
				`+DividerHtml+`
			`)
		}

		UpcomingHtml = strings.Join(UpcomingHtmlElements[:], "")
	}

	input := `
		<mjml>
			<mj-head>
				<mj-title>` + subject + `</mj-title>
				<mj-preview>` + previewText + `</mj-preview>
				<mj-style>
					.small-wrapper {
					width: 350px
					}
				</mj-style>
			</mj-head>
			<mj-body background-color="#09090B" width="500px">
				<mj-wrapper padding="32px 10px">
					<mj-section padding="16px 0px 16px 0px">
						<mj-column>
							<mj-image href="https://mailmemo.site/" alt="MailMemo logo" width="102px" height="24px" src="https://v16.appwrite.org/v1/storage/buckets/email-assets/files/668fb2360028063b4385/view?project=mail-memo" />
						</mj-column>
					</mj-section>
					<mj-section padding="0px">
						<mj-column padding="0px">
							<mj-text padding="0px" font-weight="400" font-size="32px" color="#EDEDF0" align="center">
								` + mailData.Title + `
							</mj-text>
						</mj-column>
					</mj-section>
				
					<mj-section padding="44px 0px 16px 0px" css-class="small-wrapper">
						<mj-column padding="0px">
							<mj-text font-size="18px" font-weight="400" color="#EDEDF0" padding="0px">
								Today
							</mj-text>
						</mj-column>
					</mj-section>

					` + TodayHtml + `

					<mj-section padding="32px 0px 16px 0px" css-class="small-wrapper">
						<mj-column padding="0px">
							<mj-text font-size="18px" font-weight="400" color="#EDEDF0" padding="0px">
								Upcoming Events&nbsp;
								<span style="background-color: #222224; font-size: 12px; border-radius:25px; padding: 5px 10px;">` + strconv.Itoa(totalUpcoming) + `</span>
							</mj-text>
						</mj-column>
					</mj-section>

					` + UpcomingHtml + `

					<mj-section css-class="small-wrapper">
						<mj-column>
							<mj-divider padding="24px 0px 24px 0px" border-width="1px" border-color="#222224" />
						</mj-column>
					</mj-section>
			
					<mj-section padding-bottom="16px">
						<mj-column>
							<mj-text align="center" font-weight="300" font-size="14px" color="#C3C3C6" align="center" padding="0px">
								<a href="https://app.mailmemo.site/dashboard" style="color: #C3C3C6; text-decoration: underline;">View in browser</a>
								<span style="color: #303031;">&nbsp;•&nbsp;</span>
								<a href="https://app.mailmemo.site/dashboard/settings" style="color: #C3C3C6; text-decoration: underline;">Email Preferences</a>
								<span style="color: #303031;">&nbsp;•&nbsp;</span>
								<a href="https://app.mailmemo.site/unsubscribe?email=` + userEmail + `" style="color: #C3C3C6; text-decoration: underline;">Unsubscribe</a>
							</mj-text>
						</mj-column>
					</mj-section>
			
					<mj-section padding-top="0px">
						<mj-column>
							<mj-text font-size="14px" font-weight="300" color="#FFFFFF" align="center">
							© 2024 MailMemo. All rights reserved.
							</mj-text>
						</mj-column>
					</mj-section>
				</mj-wrapper>
			</mj-body>
		</mjml>
	`

	data := map[string]interface{}{
		"subject": subject,
		"userId":  userId,
		"html":    input,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	_, executionErr := functions.CreateExecution(
		"nodeSendMailInternal",
		functions.WithCreateExecutionBody(string(jsonData)),
		functions.WithCreateExecutionAsync(true),
		functions.WithCreateExecutionMethod("POST"),
	)

	if executionErr != nil {
		Context.Error(executionErr)
		return Context.Res.Text("Error", Context.Res.WithStatusCode(500))
	}

	return Context.Res.Text("OK")
}

func getHex(color string) string {
	switch color {
	case "pink":
		return "FD366E"
	case "orange":
		return "FE9567"
	case "purple":
		return "7C67FE"
	case "mint":
		return "85DBD8"
	case "blue":
		return "68A3FE"
	case "yellow":
		return "FED367"
	}

	return "FFFFFF"
}

func arrayContains(arr []string, target string) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}
