import { databases } from '$lib/sdk';
import { user } from '$lib/stores.js';
import { redirect } from '@sveltejs/kit';

export const load = async () => {
	if (!user) {
		redirect(303, `/login`);
	}

	const calendars = await databases.listDocuments('main', 'calendars');
	const events = await databases.listDocuments('main', 'events');
	return {
		events,
		calendars
	};
};
