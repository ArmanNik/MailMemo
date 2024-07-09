package handler

import (
	"net/http"
	"time"

	"github.com/apognu/gocal"
	"github.com/open-runtimes/types-for-go/v4"
)

func Main(Context *types.Context) types.ResponseOutput {
	if Context.Req.Method != "POST" {
		return Context.Res.Text("Not Found", 404, nil)
	}

	url := Context.Req.BodyText()

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

	/*
		for _, e := range c.Events {
			// TODO: Sync event to Appwrite
		}

		gocal.Event{delayed:[]*gocal.Line{}, Uid:"3bhh4fnhlb8gb061m5f3r1a5d8@google.com", Summary:"Before", Description:"", Categories:[]string(nil), Start:time.Date(2024, time.July, 11, 0, 0, 0, 0, time.UTC), RawStart:gocal.RawDate{Params:map[string]string{"VALUE":"DATE"}, Value:"20240711"}, End:time.Date(2024, time.July, 11, 23, 59, 59, 999000000, time.UTC), RawEnd:gocal.RawDate{Params:map[string]string{"VALUE":"DATE"}, Value:"20240712"}, Duration:(*time.Duration)(nil), Stamp:time.Date(2024, time.July, 9, 15, 38, 25, 0, time.UTC), Created:time.Date(2024, time.July, 9, 13, 15, 22, 0, time.UTC), LastModified:time.Date(2024, time.July, 9, 13, 15, 22, 0, time.UTC), Location:"", Geo:(*gocal.Geo)(nil), URL:"", Status:"CONFIRMED", Organizer:(*gocal.Organizer)(nil), Attendees:[]gocal.Attendee(nil), Attachments:[]gocal.Attachment(nil), IsRecurring:false, RecurrenceID:"", RecurrenceRule:map[string]string(nil), ExcludeDates:[]time.Time(nil), Sequence:0, CustomAttributes:map[string]string(nil), Valid:true, Comment:""}
	*/

	return Context.Res.Text("OK", 200, nil)
}
