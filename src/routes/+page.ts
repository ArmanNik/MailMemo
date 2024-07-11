import { redirect } from '@sveltejs/kit';
import { base } from '$app/paths';
import { user } from '$lib/stores.js';

export async function load() {
	if (!user) {
		redirect(302, `${base}/login`);
	}
}
