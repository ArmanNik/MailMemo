import { account } from '$lib/sdk';
import { redirect } from '@sveltejs/kit';

export const load = async ({ parent }) => {
	// const { user } = await parent();
	// const prefs = await account.getPrefs();
	// if (prefs?.frequency && prefs?.interval && prefs?.time) {
	// 	redirect(303, `/dashboard/overview`);
	// } else {
	// 	redirect(303, `/dashboard/onboarding/step-1`);
	// }
};
