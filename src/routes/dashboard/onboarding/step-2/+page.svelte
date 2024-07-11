<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Select from '$lib/components/ui/select';
	import { onMount } from 'svelte';
	import { step } from '../store';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { goto } from '$app/navigation';
	import { functions } from '$lib/sdk';
	import { ExecutionMethod } from 'appwrite';

	let hour: string;
	let minute: string;
	let format: string = 'AM';
	let form: HTMLFormElement;

	onMount(() => {
		step.set(2);
	});

	async function handleNext() {
		if (!form.checkValidity()) {
			form.requestSubmit();
		}
		try {
			const cestDate = transformLocalToCEST(hour, minute, format);
			const cestHour = cestDate.getHours().toString();
			const cestMinute = cestDate.getMinutes().toString();

			const execution = await functions.createExecution(
				'api',
				JSON.stringify({
					cestHour,
					cestMinute,
					format,
					timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
				}),
				false,
				'/v1/scheduler/intervals',
				ExecutionMethod.PATCH
			);
			const isOk = execution.responseStatusCode;

			if (!isOk) {
				alert(
					execution.responseBody ? execution.responseBody : 'Unexpected error. Please try again.'
				);
				return;
			} else {
				await goto('/dashboard/onboarding/step-3');
			}
		} catch (error) {
			console.log(error);
		}
	}

	function transformLocalToCEST(hour: string, minute: string, format: string) {
		const intHour = format === 'PM' ? parseInt(hour) + 12 : parseInt(hour);
		const date = new Date();
		date.setHours(intHour);
		date.setMinutes(parseInt(minute));
		date.setSeconds(0);
		date.setMilliseconds(0);
		const cestDate = new Date(date.toLocaleString('en-US', { timeZone: 'Europe/Berlin' }));
		return cestDate;
	}

	$: if (hour && minute && format) {
		const cestDate = transformLocalToCEST(hour, minute, format);
		const cestHour = cestDate.getHours().toString();
		const cestMinute = cestDate.getMinutes().toString();
		console.log(cestHour, cestMinute);
	}
</script>

<svelte:head>
	<title>Onboarding (step 2) - MailMemo</title>
</svelte:head>

<div>
	<h1 class="font-header mt-6 text-3xl">Select the time to receive reminders</h1>
	<p class="mt-4 text-balance text-muted-foreground">
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
</div>

<div class="mt-auto flex justify-between">
	<Button variant="outline" on:click={() => goto('/dashboard/onboarding/step-1')}>Back</Button>
	<Button on:click={handleNext}>Next</Button>
</div>
