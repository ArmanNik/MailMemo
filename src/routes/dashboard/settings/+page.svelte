<script lang="ts">
	import { user } from '$lib/stores';
	import Button from '$lib/components/ui/button/button.svelte';
	import { goto, invalidate } from '$app/navigation';
	import { Dependencies } from '$lib/costants';
	import { base } from '$app/paths';
	import { account } from '$lib/sdk';

	async function logout() {
		await account.deleteSession('current');
		await invalidate(Dependencies.ACCOUNT);
		await goto(`${base}/login`);
	}
</script>

<svelte:head>
	<title>Settings - MailMemo</title>
</svelte:head>

<div class="mt-12 h-full w-full max-w-[750px] pb-16 lg:pb-0">
	<div class="flex w-full justify-between">
		<div class="flex flex-col">
			<p class="font-header text-xl">Hello</p>
			<h3 class="font-header text-2xl">
				{$user?.name}
			</h3>
		</div>
	</div>
	<img
		src="https://media3.giphy.com/media/v1.Y2lkPTc5MGI3NjExaXB1c25iZTU3cTJjYWZvdjdldWRyOWt3ZTJ0cGQ2Mm8xZm5lZzRpbCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/B9TcUZLrpj9KFD5cLw/giphy.webp"
		alt=""
	/>
	<a href="/dashboard" class="underline">Go back</a>
</div>

<div class="flex w-full justify-end">
	<Button on:click={logout}>Logout</Button>
</div>
