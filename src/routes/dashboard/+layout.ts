import { goto, invalidateAll } from '$app/navigation';
import { account } from '$lib/sdk.js';

export const load = async ({ parent, url }) => {
	// OAuth sign-in finalization
	const userId = url.searchParams.get('userId');
	const secret = url.searchParams.get('secret');
	console.log(userId, secret);
	if (secret && userId) {
		try {
			await account.createSession(userId, secret);
			await invalidateAll();
			goto('/dashboard');
			return;
		} catch (error) {
			console.warn('Sign-in flow failed');
		}
	}

	const { user } = await parent();
	return {
		user
	};
};
