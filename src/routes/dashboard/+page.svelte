<script lang="ts">
	import { user } from '$lib/stores';
	import { Badge } from '$lib/components/ui/badge';
	import { Separator } from '$lib/components/ui/separator';
	import * as Card from '$lib/components/ui/card';
	import { toLocaleTimeISO } from '$lib/utils';
	import EmptyCard from './emptyCard.svelte';
	import { colorToHex, type CalEvent } from '$lib/calendars';
	import Button from '$lib/components/ui/button/button.svelte';
	import { functions } from '$lib/sdk';
	import { ExecutionMethod } from 'appwrite';
	import { toast } from 'svelte-sonner';
	import LoaderCircle from 'lucide-svelte/icons/loader-circle';

	export let data;

	let syncButtonDisabled = false;
	let sendEmailButtonDisabled = false;

	function evaluateCornerRadius(index: number, lenght: number) {
		if (lenght === 1) {
			return 'rounded-xl';
		} else {
			//Last element
			if (index === lenght - 1) {
				return 'rounded-b-xl rounded-t-none';
			}
			//First element
			else if (index === 0) {
				return 'rounded-t-xl rounded-b-none';
			} else {
				return 'rounded-none';
			}
		}
	}

	function evaluateCornerRadiusFuture(index: number, groupedByDay: [string, CalEvent[]][]) {
		let classes = '';

		const previousGroup = groupedByDay?.[index - 1];
		const currentGroup = groupedByDay[index];
		const nextGroup = groupedByDay?.[index + 1];
		if (!previousGroup) {
			classes += ' rounded-t-xl';
		} else if (parseInt(currentGroup[0]) === parseInt(previousGroup[0]) + 1) {
			classes += ' rounded-t-none';
		}
		if (!nextGroup) {
			classes += ' rounded-b-xl';
		} else if (parseInt(currentGroup[0]) === parseInt(nextGroup[0]) - 1) {
			classes += ' rounded-b-none';
		}

		return classes;
	}

	async function syncCalendars() {
		syncButtonDisabled = true;
		try {
			const execution = await functions.createExecution(
				'syncCalendarScheduler',
				'',
				false,
				'/v1/scheduler/intervals',
				ExecutionMethod.POST
			);
			const isOk = execution.responseStatusCode;

			if (!isOk) {
				toast(
					execution.responseBody ? execution.responseBody : 'Unexpected error. Please try again.'
				);
				return;
			}
		} catch (error) {
			toast(error as string);
		} finally {
			syncButtonDisabled = false;
		}
	}

	async function sendEmail() {
		sendEmailButtonDisabled = true;
		try {
			const execution = await functions.createExecution(
				'sendMails',
				'',
				false,
				'/v1/scheduler/intervals',
				ExecutionMethod.POST
			);
			const isOk = execution.responseStatusCode === 200;

			if (!isOk) {
				toast(
					execution.responseBody ? execution.responseBody : 'Unexpected error. Please try again.'
				);
				return;
			}
		} catch (error) {
			toast(error as string);
		} finally {
			sendEmailButtonDisabled = false;
		}
	}

	$: todayEvents = data?.events?.documents?.filter((event) => {
		const eventDate = new Date(event.startAt);
		const today = new Date();
		return eventDate.getDate() === today.getDate();
	});

	$: futureEvents = data?.events?.documents?.filter((event) => {
		const eventDate = new Date(event.startAt);
		const today = new Date();
		return eventDate.getDate() > today.getDate();
	});

	//Group events that are in the same day. Use the number of days from today as the key
	$: futureEventsGrouped = futureEvents?.reduce((acc: Record<string, CalEvent[]>, event) => {
		const eventDate = new Date(event.startAt);
		const today = new Date();
		const days = eventDate.getDate() - today.getDate();
		if (!acc[days]) {
			acc[days] = [];
		}
		acc[days].push(event as CalEvent);
		return acc;
	}, {});
</script>

<svelte:head>
	<title>Dashboard - MailMemo</title>
</svelte:head>

<div class=" h-full w-full max-w-[750px] pb-16 lg:pb-0">
	<div class="flex w-full justify-between">
		<div class="flex flex-col gap-2">
			<p class="font-header text-xl font-light">Hello</p>
			<h3 class="font-header text-3xl tracking-tight">
				{$user?.name}
			</h3>
		</div>
		<a href="/dashboard/settings" aria-label="settings">
			<img src="/icons/cog.svg" alt="settings" />
		</a>
	</div>
	<div class="mt-3 flex flex-wrap gap-2">
		{#if data?.calendars?.total}
			{#each data.calendars.documents as calendar}
				<Badge variant="outline">
					<div class="h-2 w-2 rounded-full" style={`background-color: ${calendar.color}`}></div>
					<p class="ml-2">{calendar.name}</p>
				</Badge>
			{/each}
		{/if}
		<a href="/dashboard/onboarding/step-1">
			<Badge variant="outline">
				+
				<p class="ml-2">Add calendar</p>
			</Badge>
		</a>
	</div>
	<Separator class="mt-5" />
	<h1 class="font-header mt-6 text-xl">Today</h1>

	<div class="mt-4 grid gap-2">
		{#if todayEvents?.length}
			{#each todayEvents as event, i}
				{@const calendar = data.calendars.documents.find((c) => c.$id === event.calendarId)}
				<Card.Root class={`frosted ${evaluateCornerRadius(i, todayEvents?.length)} `}>
					<Card.Header class="gap-2 p-4">
						<Card.Description>
							<div class="flex h-3 items-center gap-2">
								<p>
									{toLocaleTimeISO(event.startAt, false)} - {toLocaleTimeISO(event.endAt, false)}
								</p>
								<Separator orientation="vertical" class="h-full" />

								<span class="flex items-center gap-2">
									<span
										class="h-2 w-2 rounded-full"
										style={`background-color: ${colorToHex(calendar?.color)}`}
									/>
									<p>{calendar?.name}</p>
								</span>
							</div>
						</Card.Description>

						<p class="font-header text-lg leading-4">{event.name}</p>
					</Card.Header>
				</Card.Root>
			{/each}
		{:else}
			<EmptyCard />
		{/if}
	</div>
	<h1 class="font-header mt-8 text-xl">Upcoming Events</h1>

	<div class="mt-4 grid gap-2">
		{#if futureEvents?.length}
			{@const groupedByDay = Object.entries(futureEventsGrouped)}
			{#each groupedByDay as group, groupIndex}
				<Card.Root class={`frosted ${evaluateCornerRadiusFuture(groupIndex, groupedByDay)}`}>
					<Card.Header class="p-4">
						<span class="mb-2">
							<Badge>
								In {group[0]} days
							</Badge>
						</span>
						<div class="flex flex-col gap-2">
							{#each group[1] as event, i}
								{@const calendar = data.calendars.documents.find((c) => c.$id === event.calendarId)}
								<Separator class="mt-3" />
								<Card.Description class="pt-4">
									<div class="flex h-3 items-center gap-2">
										<p>
											{toLocaleTimeISO(event.startAt, false)} - {toLocaleTimeISO(
												event.endAt,
												false
											)}
										</p>
										<Separator orientation="vertical" class="h-full" />

										<span class="flex items-center gap-2">
											<span
												class="h-2 w-2 rounded-full"
												style={`background-color: ${colorToHex(calendar?.color)}`}
											/>
											<p>{calendar?.name}</p>
										</span>
									</div>
								</Card.Description>

								<p class="font-header mt-2 text-lg leading-4">
									{event.name}
								</p>
							{/each}
						</div>
					</Card.Header>
				</Card.Root>
			{/each}
		{:else}
			<EmptyCard />
		{/if}
	</div>
</div>

<div
	class="fixed bottom-0 flex w-full max-w-[750px] justify-between gap-4 px-5 py-8 lg:relative lg:px-0"
>
	<Button
		variant="outline"
		class="w-full lg:w-auto"
		on:click={syncCalendars}
		disabled={syncButtonDisabled}>Sync calendars</Button
	>
	<Button class="w-full lg:w-auto" on:click={sendEmail} disabled={sendEmailButtonDisabled}>
		{#if sendEmailButtonDisabled}
			<LoaderCircle class="h-6 w-6" />
		{/if}
		Send email
	</Button>
</div>
