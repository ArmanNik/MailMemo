import { user } from '$lib/stores.js';
import { redirect } from '@sveltejs/kit';
import { get } from 'svelte/store';

export const load = async ({ parent }) => {
	const u = get(user);
	if (u?.prefs?.frequency && u?.prefs?.interval && u?.prefs?.time) {
		redirect(303, `/dashboard/overview`);
	} else {
		redirect(303, `/dashboard/onboarding/step-1`);
	}
};
