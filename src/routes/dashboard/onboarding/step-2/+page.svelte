<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Select from '$lib/components/ui/select';
	import { onMount } from 'svelte';
	import { step } from '../store';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { goto } from '$app/navigation';

	let hour: string;
	let minute: string;
	let period: string = 'AM';
	let form: HTMLFormElement;

	onMount(() => {
		step.set(2);
	});

	async function handleNext() {
		if (!form.checkValidity()) {
			form.requestSubmit();
		}
		try {
			await goto('/dashboard/onboarding/step-3');
		} catch (error) {
			console.log(error);
		}
	}
</script>

<div>
	<h1 class="mt-6 text-3xl font-bold">Select the time to receive reminders</h1>
	<p class="mt-4 text-balance text-muted-foreground">
		Pick the best time of day for us to send you email reminder about your upcoming events.
	</p>

	<form class="mt-20" on:submit|preventDefault bind:this={form}>
		<Label for="url">Time</Label>
		<div class="flex items-center gap-4">
			<Input id="hour" type="number" placeholder="09" required bind:value={hour} />
			<span class="text-xl font-bold">:</span>
			<Input id="hour" type="number" placeholder="00" required bind:value={minute} />
			<Select.Root
				required
				items={[{ value: 'AM' }, { value: 'PM' }]}
				onSelectedChange={(s) => {
					if (s?.value) {
						period = s.value;
					}
				}}
			>
				<Select.Trigger class="w-[73px]">
					<Select.Value placeholder="AM" />
				</Select.Trigger>
				<Select.Content>
					<Select.Item value="AM" hideCheck />
					<Select.Item value="PM" hideCheck />
				</Select.Content>
			</Select.Root>
		</div>
	</form>
</div>

<div class="mt-auto flex justify-between">
	<Button variant="outline" on:click={() => goto('/dashboard/onboarding/step-1')}>Back</Button>
	<Button on:click={handleNext}>Next</Button>
</div>
