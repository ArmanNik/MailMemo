<script lang="ts">
	import { onMount } from 'svelte';
	import '../app.css';
	import { ModeWatcher } from 'mode-watcher';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores';
	import { page } from '$app/stores';

	onMount(async () => {
		if ($user?.$id) {
			const acceptedRoutes = ['/login', '/register', '/recover'];
			if (acceptedRoutes.some((n) => $page.url.pathname.startsWith(n))) {
				await goto('/dashboard');
			}
		}
	});
</script>

<svelte:head>
	<title>MailMemo</title>
</svelte:head>

<ModeWatcher />
<slot />
