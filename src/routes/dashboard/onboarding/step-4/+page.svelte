<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import { onMount } from 'svelte';
	import { preferences, step } from '../store';
	import { Label } from '$lib/components/ui/label/index.js';
	import { goto } from '$app/navigation';
	import { account, functions } from '$lib/sdk';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import { user } from '$lib/stores';
	import * as Select from '$lib/components/ui/select';
	import { toast } from 'svelte-sonner';
	import { ExecutionMethod } from 'appwrite';

	let frequency: 'daily' | 'weekly' | 'monthly' | 'yearly' = 'daily';
	let frequencyDetails: string = '';
	let form: HTMLFormElement;

	const monthlyDetailOptions = [
		{ value: 'day1', label: 'First day' },
		{ value: 'dayLast', label: 'Last day' },
		{ value: 'dayBeforeLast', label: 'Day before last' },
		{ value: 'day7', label: 'After 7 days' },
		{ value: 'day14', label: 'After 14 days' }
	];

	const weeklyDetailOptions = [
		{ value: '0', label: 'Monday' },
		{ value: '1', label: 'Tuesday' },
		{ value: '2', label: 'Wednesday' },
		{ value: '3', label: 'Thursday' },
		{ value: '4', label: 'Friday' },
		{ value: '5', label: 'Saturday' },
		{ value: '6', label: 'Sunday' }
	];

	onMount(() => {
		step.set(4);
	});

	async function handleNext() {
		if (!form.checkValidity()) {
			form.requestSubmit();
		}
		try {
			console.log($preferences);
			await account.updatePrefs({
				...$user.prefs,
				...$preferences,
				frequency,
				frequencyDetails,
				onboarded: true
			});

			const cestDate = transformLocalToCEST(
				$preferences.hour,
				$preferences.minute,
				$preferences.format
			);
			const cestHour = cestDate.getHours();
			const cestMinute = cestDate.getMinutes();

			const execution = await functions.createExecution(
				'api',
				JSON.stringify({
					hours: cestHour,
					minutes: cestMinute,
					format: $preferences.format.toLowerCase(),
					frequency,
					frequencyDetails
				}),
				false,
				'/v1/scheduler/intervals',
				ExecutionMethod.PATCH
			);
			const isOk = execution.responseStatusCode === 200;

			if (!isOk) {
				toast(
					execution.responseBody ? execution.responseBody : 'Unexpected error. Please try again.'
				);
				return;
			} else {
				await goto('/dashboard');
			}
		} catch (error) {
			toast(error as string);
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
</script>

<svelte:head>
	<title>Onboarding: set frequency - MailMemo</title>
</svelte:head>

<div>
	<h1 class="mt-6 max-w-[80%] font-header text-3xl tracking-tight">
		Set the frequency of receiving reminders
	</h1>
	<p class="mt-4 text-muted-foreground">Decide how often you want to receive email reminders.</p>

	<form class="mt-20" on:submit|preventDefault bind:this={form}>
		<RadioGroup.Root bind:value={frequency}>
			<div class="flex flex-col gap-4">
				<div class="flex min-h-10 items-center space-x-2 px-1">
					<RadioGroup.Item value="daily" id="daily" />
					<Label for="daily">Daily</Label>
				</div>
				<div class="flex min-h-10 items-center gap-5">
					<div class="flex items-center space-x-2 px-1">
						<RadioGroup.Item value="weekly" id="weekly" />
						<Label for="weekly">Weekly</Label>
					</div>
					{#if frequency === 'weekly'}
						<Select.Root
							items={weeklyDetailOptions}
							onSelectedChange={(s) => {
								if (s?.value) {
									frequencyDetails = s.value;
								}
							}}
						>
							<Select.Trigger class=" w-full lg:w-[180px]">
								<Select.Value placeholder="Day in week" />
							</Select.Trigger>
							<Select.Content>
								{#each weeklyDetailOptions as option}
									<Select.Item value={option.value} hideCheck>{option.label}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					{/if}
				</div>
				<div class="flex min-h-10 items-center gap-5">
					<div class="flex items-center space-x-2 px-1">
						<RadioGroup.Item value="monthly" id="monthly" />
						<Label for="monthly">Monthly</Label>
					</div>
					{#if frequency === 'monthly'}
						<Select.Root
							items={monthlyDetailOptions}
							onSelectedChange={(s) => {
								if (s?.value) {
									frequencyDetails = s.value;
								}
							}}
						>
							<Select.Trigger class="w-full lg:w-[180px]">
								<Select.Value placeholder="Day in month" />
							</Select.Trigger>
							<Select.Content>
								{#each monthlyDetailOptions as option}
									<Select.Item value={option.value} hideCheck>{option.label}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					{/if}
				</div>
			</div>
		</RadioGroup.Root>
	</form>
</div>

<div class="mt-auto flex justify-between">
	<Button variant="outline" on:click={() => goto('/dashboard/onboarding/step-3')}>Back</Button>
	<Button on:click={handleNext}>Next</Button>
</div>
