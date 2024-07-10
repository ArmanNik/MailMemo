package handler

import (
	"openruntimes/handler/services"
	"os"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/open-runtimes/types-for-go/v4"
)

func Main(Context *types.Context) types.ResponseOutput {
	appwriteClient := client.NewClient()
	appwriteClient.SetEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT"))
	appwriteClient.SetProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID"))
	appwriteClient.SetKey(Context.Req.Headers["x-appwrite-key"])

	action := Context.Req.Method + " " + Context.Req.Path

	if action == "PATCH /v1/scheduler/intervals" {
		return services.UpdateSchedulerInterval(Context, appwriteClient)
	}

	return Context.Res.Text("Not Found", 404, nil)
}
