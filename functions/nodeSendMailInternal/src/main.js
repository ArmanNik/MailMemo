import mjml2html from 'mjml';
import { Client, Messaging, ID } from 'node-appwrite';

export default async ({ req, res, log }) => {
  const bodyText = req.bodyText;

  const json = JSON.parse(bodyText ?? '{}');

  const client = new Client()
    .setEndpoint(process.env.APPWRITE_FUNCTION_API_ENDPOINT)
    .setProject(process.env.APPWRITE_FUNCTION_PROJECT_ID)
    .setKey(req.headers['x-appwrite-key']);

  const messaging = new Messaging(client);

  const htmlOutput = mjml2html(json.html, {
    beautify: false,
    keepComments: false,
    validationLevel: 'skip'
  });

  const message = await messaging.createEmail(
    ID.unique(),
    json.subject,
    htmlOutput.html,
    undefined,
    [json.userId],
    undefined,
    undefined,
    undefined,
    undefined,
    undefined,
    true,
    undefined);

  log(message.$id);

  return res.send("OK");
};
