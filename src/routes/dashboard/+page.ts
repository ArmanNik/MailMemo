import { Dependencies } from '$lib/costants.js';
import { databases } from '$lib/sdk';
import { user } from '$lib/stores.js';
import { redirect } from '@sveltejs/kit';

export const load = async ({ depends }) => {
	if (!user) {
		redirect(303, `/login`);
	}
	depends(Dependencies.DOCUMENTS);
	const calendars = await databases.listDocuments('main', 'calendars');
	const events = await databases.listDocuments('main', 'events');
	return {
		events,
		calendars
	};
};
