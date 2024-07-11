<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(async () => {
		if ($page?.status === 404) {
			await goto('/dashboard');
		}
	});
</script>

<h1>{$page?.error?.message}</h1>
<a href="/">Back to the console</a>
