import { Client, Account, Databases, OAuthProvider } from 'appwrite';

const client = new Client();

client.setEndpoint('https://v16.appwrite.org/v1').setProject('mail-memo');

export const account = new Account(client);
account.createOAuth2Session(
	OAuthProvider.Github, // provider
	'https://example.com/success', // redirect here on success
	'https://example.com/failed', // redirect here on failure
	['repo', 'user'] // scopes (optional)
);
export const databases = new Databases(client);
