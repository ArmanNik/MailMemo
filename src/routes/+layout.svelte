<script lang="ts">
	import { onMount } from 'svelte';
	import '../app.css';
	import { ModeWatcher } from 'mode-watcher';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores';

	onMount(async () => {
		if ($user?.$id) {
			await goto('/dashboard');
		}
	});
</script>

<ModeWatcher />
<slot />
