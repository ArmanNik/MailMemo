package handler

import (
	"openruntimes/handler/services"
	"os"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/open-runtimes/types-for-go/v4"
)

func Main(context types.Context) types.ResponseOutput {
	client := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_API_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_PROJECT_ID")),
		appwrite.WithKey(context.Req.Headers["x-appwrite-key"]),
	)

	action := context.Req.Method + " " + context.Req.Path
	switch a := action; a {
	case "PATCH /v1/scheduler/intervals":
		return services.UpdateSchedulerInterval(context, client)
	case "POST /v1/calendars":
		return services.CreateCalendar(context, client)
	case "DELETE /v1/subscriptions":
		return services.DeleteSubscription(context, client)
	default:
		return context.Res.Text("Not Found", context.Res.WithStatusCode(404))
	}
}
