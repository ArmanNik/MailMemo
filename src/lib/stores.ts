import { page } from '$app/stores';
import type { Models } from 'appwrite';
import { derived } from 'svelte/store';

export const user = derived(page, ($page) => {
	return $page.data.user as Models.User<Record<string, string>>;
});
