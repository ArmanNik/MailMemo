<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { page } from '$app/stores';
	import { Dependencies } from '$lib/costants';
	import { client } from '$lib/sdk';
	import { user } from '$lib/stores';
	import type { Models } from 'appwrite';
	import { onMount } from 'svelte';

	onMount(() => {
		if (!$user?.prefs?.onboarded && !$page.url.pathname.includes('/onboarding')) {
			if ($user?.prefs?.firstCal) {
				goto('/dashboard/onboarding/step-2');
			} else {
				goto('/dashboard/onboarding/step-1');
			}
			return;
		}
		let previousStatus: string | null = null;

		return client.subscribe<Models.Document>('documents', (message) => {
			if (previousStatus === message.payload.status) {
				return;
			}
			previousStatus = message.payload.status;
			if (message.events.includes('databases.*.collections.*.create')) {
				invalidate(Dependencies.DOCUMENTS);

				return;
			}
			if (message.events.includes('databases.*.collections.*.update')) {
				invalidate(Dependencies.DOCUMENTS);

				return;
			}
			if (message.events.includes('databases.*.collections.*.delete')) {
				invalidate(Dependencies.DOCUMENTS);

				return;
			}
		});
	});
</script>

<div
	class="mm-container flex h-full min-h-[100vh] grow flex-col items-center justify-between px-5 py-8"
	class:mm-container={!$page.url.pathname.includes('/onboarding')}
>
	<slot />
</div>
