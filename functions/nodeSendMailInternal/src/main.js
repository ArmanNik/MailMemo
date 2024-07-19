import FormData from 'form-data';
import Mailgun from 'mailgun.js';
import mjml2html from 'mjml';
import { Client, Users } from 'node-appwrite';

export default async ({ req, res, log }) => {
  const bodyText = req.bodyText;

  const json = JSON.parse(bodyText ?? '{}');

  const client = new Client()
    .setEndpoint(process.env.APPWRITE_FUNCTION_API_ENDPOINT)
    .setProject(process.env.APPWRITE_FUNCTION_PROJECT_ID)
    .setKey(req.headers['x-appwrite-key']);

  const users = new Users(client);

  const user = await users.get(json.userId);

  const userEmail = user.email;

  const htmlOutput = mjml2html(json.html, {
    beautify: false,
    keepComments: false,
    validationLevel: 'skip'
  });

  const mailgun = new Mailgun(FormData);
  const mg = mailgun.client({
    url: 'https://api.eu.mailgun.net', // To use EU domains
    username: 'api',
    key: process.env.MAILGUN_API_KEY
  });

  const data = {
    from: "noreply@mg.almostapps.eu",
    to: userEmail,
    subject: json.subject,
    html: htmlOutput.html
  };

  const response = await mg.messages.create("mg.almostapps.eu", data);

  log(response.id);

  return res.send("OK");
};
