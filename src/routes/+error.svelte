<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	onMount(() => {
		console.log('test');
	});
</script>

<h1>{$page?.error?.message}</h1>
<a href="/">Back to the console</a>
