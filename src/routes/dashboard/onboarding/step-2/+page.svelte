<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Select from '$lib/components/ui/select';
	import { onMount } from 'svelte';
	import { preferences, step } from '../store';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { goto } from '$app/navigation';
	import { functions } from '$lib/sdk';
	import { ExecutionMethod } from 'appwrite';
	import { toast } from 'svelte-sonner';

	let hour: string = '09';
	let minute: string = '00';
	let format: string = 'AM';
	let form: HTMLFormElement;

	onMount(() => {
		step.set(2);
	});

	async function handleNext() {
		if (!form.checkValidity()) {
			form.requestSubmit();
		}
		$preferences = {
			...$preferences,
			hour,
			minute,
			format,
			timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
		};
		await goto('/dashboard/onboarding/step-3');
	}

	$: if (parseInt(hour) > 12) {
		hour = '12';
	}
	$: if (parseInt(minute) > 59) {
		minute = '59';
	}
	$: if (parseInt(hour) < 0) {
		hour = '0';
	}
	$: if (parseInt(minute) < 0) {
		minute = '0';
	}
</script>

<svelte:head>
	<title>Onboarding: select time - MailMemo</title>
</svelte:head>

<div>
	<h1 class="mt-6 max-w-[80%] font-header text-3xl tracking-tight">
		Select the time to receive reminders
	</h1>
	<p class="mt-4 text-muted-foreground">
		Pick the best time of day for us to send you email reminder about your upcoming events.
	</p>

	<form class="mt-20" on:submit|preventDefault bind:this={form}>
		<Label for="url">Time</Label>
		<div class="flex items-center gap-4">
			<Input
				id="hour"
				type="number"
				placeholder="Hour"
				min="0"
				max="12"
				required
				bind:value={hour}
			/>
			<span class="text-xl font-bold">:</span>
			<Input
				id="hour"
				type="number"
				placeholder="Minutes"
				min="0"
				max="59"
				required
				bind:value={minute}
			/>
			<Select.Root
				required
				items={[{ value: 'AM' }, { value: 'PM' }]}
				onSelectedChange={(s) => {
					if (s?.value) {
						format = s.value;
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
	<p class="mt-2 text-muted-foreground">
		Your reminder email will be sent at {format === 'PM'
			? parseInt(hour) + 12
			: parseInt(hour)}:{minute}.
	</p>
</div>

<div class="mt-auto flex justify-end">
	<!-- <Button variant="outline" on:click={() => goto('/dashboard/onboarding/step-1')}>Back</Button> -->
	<Button on:click={handleNext}>Next</Button>
</div>
