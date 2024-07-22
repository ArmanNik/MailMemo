package handler

import (
	"openruntimes/handler/services"
	"os"

	"github.com/appwrite/sdk-for-go/client"
	openruntimes "github.com/open-runtimes/types-for-go/v4"
)

func Main(Context openruntimes.Context) openruntimes.Response {
	appwriteClient := client.NewClient()
	appwriteClient.SetEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT"))
	appwriteClient.SetProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID"))
	appwriteClient.SetKey(Context.Req.Headers["x-appwrite-key"])

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
