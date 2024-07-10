<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { user } from '$lib/stores';
	import { onMount } from 'svelte';

	onMount(async () => {
		if (!$user?.prefs?.onboarded && !$page.url.pathname.includes('/onboarding')) {
			await goto('/dashboard/onboarding/step-1');
		}
	});
</script>

<div class="flex h-full min-h-[100vh] grow flex-col items-center justify-between px-5 py-8">
	<slot />
</div>
