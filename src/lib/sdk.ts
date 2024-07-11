import { Client, Account, Databases, Functions } from 'appwrite';

export const client = new Client();

client.setEndpoint('https://v16.appwrite.org/v1').setProject('mail-memo');

export const account = new Account(client);
export const databases = new Databases(client);
export const functions = new Functions(client);
