import { account } from '$lib/sdk';
import { redirect } from '@sveltejs/kit';

export const ssr = false;

export const load = async ({ depends, url }) => {
	// depends(Dependencies.ACCOUNT);

	try {
		const user = await account.get();
		return {
			user
		};
	} catch (error) {
		const acceptedRoutes = ['/login', '/register', '/recover'];
		console.log(error);
		if (!acceptedRoutes.some((n) => url.pathname.startsWith(n))) {
			redirect(303, `/login`);
		}
		return {};
	}
};
