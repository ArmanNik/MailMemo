import { user } from '$lib/stores.js';
import { redirect } from '@sveltejs/kit';

export const load = async () => {
	if (!user) {
		redirect(303, `/login`);
	}
};
