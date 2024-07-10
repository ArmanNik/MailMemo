<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import { onMount } from 'svelte';
	import { step } from '../store';
	import { Label } from '$lib/components/ui/label/index.js';
	import { goto } from '$app/navigation';
	import { account } from '$lib/sdk';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import { user } from '$lib/stores';

	let frequency: 'daily' | 'weekly' | 'monthly' | 'yearly' = 'daily';
	let form: HTMLFormElement;

	onMount(() => {
		step.set(4);
	});

	async function handleNext() {
		if (!form.checkValidity()) {
			form.requestSubmit();
		}
		try {
			await account.updatePrefs({
				...$user.prefs,
				frequency
			});
			await goto('/dashboard/onboarding/overview');
		} catch (error) {
			console.log(error);
		}
	}
</script>

<div>
	<h1 class="mt-6 text-3xl font-bold">Set the frequency of receiving reminders</h1>
	<p class="mt-4 text-balance text-muted-foreground">
		Decide how often you want to receive email reminders.
	</p>

	<form class="mt-20" on:submit|preventDefault bind:this={form}>
		<RadioGroup.Root bind:value={frequency}>
			<div class="flex flex-col gap-5">
				<div class="flex items-center space-x-2 px-1 py-2">
					<RadioGroup.Item value="daily" id="daily" />
					<Label for="daily">Daily</Label>
				</div>
				<div class="flex items-center space-x-2 px-1 py-2">
					<RadioGroup.Item value="weekly" id="weekly" />
					<Label for="weekly">Weekly</Label>
				</div>
				<div class="flex items-center space-x-2 px-1 py-2">
					<RadioGroup.Item value="monthly" id="monthly" />
					<Label for="monthly">Monthly</Label>
				</div>
			</div>
		</RadioGroup.Root>
	</form>
</div>

<div class="mt-auto flex justify-between">
	<Button variant="outline" on:click={() => goto('/dashboard/onboarding/step-3')}>Back</Button>
	<Button on:click={handleNext}>Next</Button>
</div>
