package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/apognu/gocal"
	"github.com/open-runtimes/types-for-go/v4"
)

// This is your Appwrite function
// It's executed each time we get a request
func Main(Context *types.Context) types.ResponseOutput {
	if Context.Req.Method != "POST" {
		return Context.Res.Text("Not Found", 404, nil)
	}

	url := Context.Req.BodyText()

	if os.Getenv("APPWRITE_ENV") == "development" {
		url = "https://calendar.google.com/calendar/ical/e251560319ad845251b578e5a88962f3c20cf07a6900a462f75cb42c7dc898ca%40group.calendar.google.com/private-67520c9a35acde3e6ab144529aa7f389/basic.ics"
	}

	calResp, err := http.Get(url)
	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", 500, nil)
	}

	defer calResp.Body.Close()

	start, end := time.Now(), time.Now().Add(12*30*24*time.Hour)

	c := gocal.NewParser(calResp.Body)
	c.Start, c.End = &start, &end
	c.Parse()

	Context.Log("Reading")

	for _, e := range c.Events {
		Context.Log(e.Summary)
	}

	Context.Log("End of it")

	return Context.Res.Text("OK", 200, nil)
}
