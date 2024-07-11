<script lang="ts">
	import { onMount } from 'svelte';
	import '../app.css';
	import { ModeWatcher } from 'mode-watcher';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores';
	import { page } from '$app/stores';
	import { Toaster } from '$lib/components/ui/sonner';

	onMount(async () => {
		if ($user?.$id) {
			const skippedRoutes = ['/login', '/register', '/recover'];
			if (
				skippedRoutes.some((n) => $page.url.pathname.startsWith(n)) ||
				$page.url.pathname === '/'
			) {
				await goto('/dashboard');
			}
		}
	});
</script>

<ModeWatcher defaultMode="dark" />
<Toaster position="top-right" />

<slot />
