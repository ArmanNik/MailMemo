<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import { onMount } from 'svelte';
	import { preferences, step } from '../store';
	import { Label } from '$lib/components/ui/label/index.js';
	import { goto } from '$app/navigation';
	import { account } from '$lib/sdk';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import { user } from '$lib/stores';

	let period: 'today' | 'week' | 'month' | 'year' = 'year';
	let form: HTMLFormElement;

	onMount(() => {
		step.set(3);
	});

	async function handleNext() {
		if (!form.checkValidity()) {
			form.requestSubmit();
		}
		preferences.set({ period });
		await goto('/dashboard/onboarding/step-4');

		// try {
		// 	await account.updatePrefs({
		// 		...$user.prefs,
		// 		period
		// 	});
		// } catch (error) {
		// 	console.log(error);
		// }
	}
</script>

<svelte:head>
	<title>Onboarding: set period - MailMemo</title>
</svelte:head>

<div>
	<h1 class="mt-6 max-w-[80%] font-header text-3xl tracking-tight">
		Set period for event reminders
	</h1>
	<p class="mt-4 text-muted-foreground">
		Decide how far in advance you'd like to be reminded of your events.
	</p>

	<form class="mt-20" on:submit|preventDefault bind:this={form}>
		<RadioGroup.Root bind:value={period}>
			<div class="flex flex-col gap-5">
				<div class="flex items-center space-x-2 px-1 py-2">
					<RadioGroup.Item value="today" id="today" />
					<Label for="today">Today</Label>
				</div>
				<div class="flex items-center space-x-2 px-1 py-2">
					<RadioGroup.Item value="week" id="week" />
					<Label for="week">This week</Label>
				</div>
				<div class="flex items-center space-x-2 px-1 py-2">
					<RadioGroup.Item value="month" id="month" />
					<Label for="month">This month</Label>
				</div>
				<div class="flex items-center space-x-2 px-1 py-2">
					<RadioGroup.Item value="year" id="year" />
					<Label for="year">This year</Label>
				</div>
			</div>
		</RadioGroup.Root>
	</form>
</div>

<div class="mt-auto flex justify-between">
	<Button variant="outline" on:click={() => goto('/dashboard/onboarding/step-2')}>Back</Button>
	<Button on:click={handleNext}>Next</Button>
</div>
