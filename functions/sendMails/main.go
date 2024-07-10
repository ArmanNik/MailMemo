package handler

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/Boostport/mjml-go"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/messaging"
	"github.com/open-runtimes/types-for-go/v4"
)

func Main(Context *types.Context) types.ResponseOutput {
	if Context.Req.Method != "POST" {
		return Context.Res.Text("Not Found", 404, nil)
	}

	appwriteClient := client.NewClient()
	appwriteClient.SetEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT"))
	appwriteClient.SetProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID"))
	appwriteClient.SetKey(Context.Req.Headers["x-appwrite-key"])

	appwriteMessaging := messaging.NewMessaging(appwriteClient)

	subject := "MJML E-mail"
	input := `
	<mjml>
	<mj-body>
	  <mj-wrapper border="1px solid #EDEDF0" padding="50px 30px" background-color="#FAFAFB">
		<mj-section>
		  <mj-column>
			<mj-text font-size="20px" color="#56565C" font-family="helvetica" padding-left="0px">
			  July 9, 2024
			</mj-text>
			<mj-text font-size="30px" color="#19191C" font-family="helvetica" padding-left="0px">
						  Upcoming Events
			</mj-text>
		  </mj-column>
		</mj-section>

		<mj-section>
			<mj-column>
				<mj-text font-size="50px" color="#56565C" font-family="helvetica">•</mj-text>
			</mj-column>
			<mj-column>
				<mj-text font-size="12px" color="#2D2D31" font-family="helvetica" padding-top="0px">
					Submit final project  adioadsno noiasd noisdaoin aonisadf oni asfdionafnd sioreport
				</mj-text>
			</mj-column>
		</mj-section>
		

		<mj-section padding="5px 0">  </mj-section>
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
			Context.Req.Headers["x-appwrite-user-id"],
		}),
	)

	if err != nil {
		Context.Error(err)
		return Context.Res.Text("Error", 500, nil)
	}

	Context.Log("Message ID: " + message.Id)

	return Context.Res.Text("OK", 200, nil)
}
