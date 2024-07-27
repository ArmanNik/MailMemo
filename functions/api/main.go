package handler

import (
	"openruntimes/handler/services"
	"os"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

func Main(Context openruntimes.Context) openruntimes.Response {
	appwriteClient := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID")),
		appwrite.WithKey(Context.Req.Headers["x-appwrite-key"]),
	)

	action := Context.Req.Method + " " + Context.Req.Path
	switch a := action; a {
	case "PATCH /v1/scheduler/intervals":
		return services.UpdateSchedulerInterval(Context, appwriteClient)
	case "POST /v1/calendars":
		return services.CreateCalendar(Context, appwriteClient)
	case "DELETE /v1/subscriptions":
		return services.DeleteSubscription(Context, appwriteClient)
	default:
		return Context.Res.Text("Not Found", Context.Res.WithStatusCode(404))
	}
}
